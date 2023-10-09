package controller

import (
	"encoding/json"
	"math/rand"

	"github.com/farseeingnorthwest/battleground.go/functional"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/gofiber/fiber/v2"
)

type BattleController struct {
	CharacterRepo CharacterRepository
	skillRepo     SkillRepository
}

func NewBattleController(characterRepo CharacterRepository, skillRepo SkillRepository) BattleController {
	return BattleController{characterRepo, skillRepo}
}

func (c BattleController) Mount(router fiber.Router) {
	router.Post("/battles", c.CreateBattle)
}

func (c BattleController) CreateBattle(fc *fiber.Ctx) error {
	form := struct {
		Seed   int64
		Left   map[int]int
		Right  map[int]int
		Ground []int
	}{}
	if err := fc.BodyParser(&form); err != nil {
		return err
	}

	skills, err := c.getSkills()
	if err != nil {
		return err
	}

	left, err := c.getWarriors(form.Left, battlefield.Left, skills)
	if err != nil {
		return err
	}
	right, err := c.getWarriors(form.Right, battlefield.Right, skills)
	if err != nil {
		return err
	}
	ob := newObserver()
	f := battlefield.NewBattleField(
		rand.New(rand.NewSource(form.Seed)),
		append(left, right...),
		append(
			functional.MapSlice(
				func(id int) battlefield.Reactor {
					return skills[id].Reactor.Spawn()
				},
				form.Ground,
			),
			ob,
			newObserverComplete(ob),
		)...,
	)
	f.Run()

	return fc.JSON(ob)
}

func (c BattleController) getSkills() (map[int]storage.Skill, error) {
	skills, err := c.skillRepo.FindEx()
	if err != nil {
		return nil, err
	}

	return functional.Tabulate[int, storage.Skill](bySkillID(skills)), nil
}

func (c BattleController) getWarriors(m map[int]int, side battlefield.Side, skills map[int]storage.Skill) ([]battlefield.Warrior, error) {
	charSlice, err := c.CharacterRepo.Find(functional.Values(m)...)
	if err != nil {
		return nil, err
	}

	warriors := make([]battlefield.Warrior, 0, len(m))
	chars := functional.Tabulate[int, storage.Character](byCharacterID(charSlice))
	for p, id := range m {
		warriors = append(warriors, battlefield.NewMyWarrior(
			battlefield.MyBaseline{
				Damage:       chars[id].Damage,
				CriticalOdds: chars[id].CriticalOdds,
				CriticalLoss: chars[id].CriticalLoss,
				Defense:      chars[id].Defense,
				Health:       chars[id].Health,
				Speed:        chars[id].Speed,
			},
			side,
			p,
			battlefield.WarriorSkills(functional.MapSlice(
				func(meta storage.SkillMeta) battlefield.Reactor {
					return skills[meta.ID].Reactor.Spawn()
				},
				functional.Values(chars[id].Skills),
			)...),
		))
	}

	return warriors, nil
}

type observer struct {
	battlefield.TagSet
	rounds []*round
	winner battlefield.Warrior
}

func newObserver() *observer {
	return &observer{battlefield.NewTagSet(battlefield.Priority(1000000)), nil, nil}
}

func (o *observer) top() *round {
	return o.rounds[len(o.rounds)-1]
}

func (o *observer) React(signal battlefield.Signal, ec battlefield.EvaluationContext) {
	switch signal := signal.(type) {
	case *battlefield.BattleStartSignal:
		o.rounds = append(o.rounds, newRound(1, ec))
	case *battlefield.RoundStartSignal:
		o.rounds = append(o.rounds, newRound(0, ec))
	case *battlefield.RoundEndSignal:
		o.top().setCurrent(2)
	case *battlefield.LifecycleSignal:
		o.top().appendSignal(signal)
	case *battlefield.PostActionSignal:
		o.top().appendAction(signal.Action())
		if _, source, _ := signal.Action().Script().Source(); source != nil {
			o.winner = source.(battlefield.Warrior)
		}
	}
}

func (o *observer) Active() bool {
	return true
}

func (o *observer) MarshalJSON() ([]byte, error) {
	start := o.rounds[0].stages[1]
	if start == nil {
		start = act{}
	}

	return json.Marshal(map[string]any{
		"winner":   o.winner.Side().String(),
		"profiles": o.rounds[0].profiles,
		"start":    start,
		"rounds":   o.rounds[1:],
	})
}

type observerComplete struct {
	battlefield.TagSet
	*observer
}

func newObserverComplete(ob *observer) *observerComplete {
	return &observerComplete{battlefield.NewTagSet(battlefield.Priority(-1000000)), ob}
}

func (c *observerComplete) React(signal battlefield.Signal, _ battlefield.EvaluationContext) {
	if _, ok := signal.(*battlefield.RoundStartSignal); ok {
		c.top().setCurrent(1)
	}
}

type round struct {
	profiles []*profile
	stages   [3]act
	current  int
}

func newRound(current int, ec battlefield.EvaluationContext) *round {
	return &round{createProfiles(ec), [3]act{}, current}
}

func (r *round) setCurrent(current int) {
	r.current = current
}

func (r *round) appendAction(action battlefield.Action) {
	r.stages[r.current] = append(r.stages[r.current], sentence{action: newActionView(action)})
}

func (r *round) appendSignal(signal *battlefield.LifecycleSignal) {
	r.stages[r.current] = append(r.stages[r.current], sentence{signal: newSignalView(signal)})
}

func (r *round) MarshalJSON() ([]byte, error) {
	v := map[string]any{
		"profiles": r.profiles,
	}
	for i, stage := range r.stages {
		if len(stage) > 0 {
			v[[]string{"start", "main", "end"}[i]] = stage
		}
	}

	return json.Marshal(v)
}

type act []sentence

type sentence struct {
	action *actionView
	signal *signalView
}

func (s sentence) MarshalJSON() ([]byte, error) {
	if s.action != nil {
		return json.Marshal(s.action)
	}

	return json.Marshal(s.signal)
}

type signalView struct {
	signal battlefield.Signal
}

func newSignalView(signal battlefield.Signal) *signalView {
	return &signalView{signal}
}

func (v signalView) MarshalJSON() ([]byte, error) {
	switch signal := v.signal.(type) {
	case *battlefield.LifecycleSignal:
		scripter, reactor := signal.Source()
		v := map[string]any{
			"id":      signal.ID(),
			"name":    signal.Name(),
			"parent":  newSignalView(signal.Parent()),
			"reactor": newReactorView(reactor),
		}
		if w, ok := scripter.(battlefield.Warrior); ok {
			v["warrior"] = newWarriorView(w)
		}
		if signal.Lifecycle() != nil {
			v["lifecycle"] = (*lifecycleView)(signal.Lifecycle())
		}
		if signal.Affairs() != 0 {
			v["affairs"] = affairsView(signal.Affairs())
		}

		return json.Marshal(v)
	case battlefield.ActionSignal:
		return json.Marshal(map[string]any{
			"id":     signal.ID(),
			"name":   signal.Name(),
			"action": signal.Action().ID(),
		})
	default:
		return json.Marshal(map[string]any{
			"id":   signal.ID(),
			"name": signal.Name(),
		})
	}
}

type warriorView struct {
	warrior battlefield.Warrior
}

func newWarriorView(warrior battlefield.Warrior) warriorView {
	return warriorView{warrior}
}

func (v warriorView) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"side":     v.warrior.Side().String(),
		"position": v.warrior.Position(),
	})
}

type reactorView struct {
	reactor battlefield.Reactor
}

func newReactorView(reactor battlefield.Reactor) reactorView {
	return reactorView{reactor}
}

func (v reactorView) MarshalJSON() ([]byte, error) {
	if label, ok := battlefield.QueryTag[battlefield.Label](v.reactor); ok {
		return json.Marshal(string(label))
	}

	return json.Marshal("Unknown")
}

type sourceView struct {
	signal   battlefield.Signal
	scripter any
	reactor  battlefield.Reactor
}

func newSourceView(signal battlefield.Signal, scripter any, reactor battlefield.Reactor) sourceView {
	return sourceView{signal, scripter, reactor}
}

func (s sourceView) MarshalJSON() ([]byte, error) {
	v := map[string]any{
		"signal":  newSignalView(s.signal),
		"reactor": newReactorView(s.reactor),
	}
	if w, ok := s.scripter.(battlefield.Warrior); ok {
		v["warrior"] = newWarriorView(w)
	}

	return json.Marshal(v)
}

type lifecycleView battlefield.Lifecycle

func (lc *lifecycleView) MarshalJSON() ([]byte, error) {
	v := make(map[string]int)
	if lc.Leading.Ok() {
		v["leading"] = lc.Leading.Value()
	}
	if lc.Cooling.Ok() {
		v["cooling"] = lc.Cooling.Value().Current
	}
	if lc.Capacity.Ok() {
		v["capacity"] = lc.Capacity.Value()
	}

	return json.Marshal(v)
}

type affairsView battlefield.LifecycleAffairs

func (a affairsView) MarshalJSON() ([]byte, error) {
	var v []string
	if a&affairsView(battlefield.LifecycleTrigger) != 0 {
		v = append(v, "Trigger")
	}
	if a&affairsView(battlefield.LifecycleOverflow) != 0 {
		v = append(v, "Overflow")
	}

	return json.Marshal(v)
}

type actionView struct {
	action    battlefield.Action
	health    map[battlefield.Warrior]healthView
	lifecycle map[battlefield.Reactor]lifecycleView
}

func newActionView(action battlefield.Action) *actionView {
	v := &actionView{action: action}
	switch verb := action.Verb().(type) {
	case *battlefield.Attack:
		v.health = functional.MapValues2(func(_ int, w battlefield.Warrior) healthView {
			return healthView(w.Health())
		}, verb.Loss())
	case *battlefield.Heal:
		v.health = functional.MapValues2(func(_ int, w battlefield.Warrior) healthView {
			return healthView(w.Health())
		}, verb.Rise())
	case *battlefield.Buff:
		v.lifecycle = make(map[battlefield.Reactor]lifecycleView)
		for _, r := range verb.Provision() {
			v.lifecycle[r] = lifecycleView(r.(*battlefield.FatReactor).Lifecycle())
		}
		for _, r := range verb.Overflow() {
			v.lifecycle[r] = lifecycleView(r.(*battlefield.FatReactor).Lifecycle())
		}
	case *battlefield.Purge:
		v.lifecycle = make(map[battlefield.Reactor]lifecycleView)
		for _, rs := range verb.Recycles() {
			for _, r := range rs {
				v.lifecycle[r] = lifecycleView(r.(*battlefield.FatReactor).Lifecycle())
			}
		}
	}

	return v
}

func (a actionView) createEvolution(warrior battlefield.Warrior, value int) evolution {
	return evolution{
		Warrior: newWarriorView(warrior),
		Health:  a.health[warrior],
		Value:   value,
	}
}

func (a actionView) createProvision(warrior battlefield.Warrior, reactor battlefield.Reactor) provision {
	return provision{
		Warrior:   newWarriorView(warrior),
		Lifecycle: lifecycleView(a.lifecycle[reactor]),
	}
}

func (a actionView) MarshalJSON() ([]byte, error) {
	v := map[string]any{
		"id":      a.action.ID(),
		"source":  newSourceView(a.action.Script().Source()),
		"targets": functional.MapSlice(newWarriorView, a.action.Targets()),
	}
	if len(a.action.FalseTargets()) > 0 {
		v["false_targets"] = functional.MapSlice(newWarriorView, a.action.FalseTargets())
	}
	if len(a.action.ImmuneTargets()) > 0 {
		v["immune_targets"] = functional.MapSlice(newWarriorView, functional.Keys(a.action.ImmuneTargets()))
	}

	switch verb := a.action.Verb().(type) {
	case *battlefield.Attack:
		v["verb"] = map[string]any{
			"_verb":    "attack",
			"critical": verb.Critical(),
			"losses":   functional.MapKVs(a.createEvolution, verb.Loss()),
		}

	case *battlefield.Heal:
		v["verb"] = map[string]any{
			"_verb": "heal",
			"rises": functional.MapKVs(a.createEvolution, verb.Rise()),
		}

	case *battlefield.Buff:
		v["verb"] = map[string]any{
			"_verb":      "buff",
			"reactor":    newReactorView(verb.Reactor()),
			"provisions": functional.MapKVs(a.createProvision, verb.Provision()),
			"overflows":  functional.MapKVs(a.createProvision, verb.Overflow()),
		}

	case *battlefield.Purge:
		v["verb"] = map[string]any{
			"_verb": "purge",
			"recycles": functional.MapKVs(func(warrior battlefield.Warrior, reactors []battlefield.Reactor) recycle {
				return recycle{
					Warrior: newWarriorView(warrior),
					Reactors: functional.MapSlice(func(reactor battlefield.Reactor) reactorLifecycleView {
						return reactorLifecycleView{
							newReactorView(reactor),
							lifecycleView(a.lifecycle[reactor]),
						}
					}, reactors),
				}
			}, verb.Recycles()),
		}
	}

	return json.Marshal(v)
}

type evolution struct {
	Warrior warriorView `json:"warrior"`
	Health  healthView  `json:"health"`
	Value   int         `json:"value"`
}

type provision struct {
	Warrior   warriorView   `json:"warrior"`
	Lifecycle lifecycleView `json:"lifecycle"`
}

type recycle struct {
	Warrior  warriorView            `json:"warrior"`
	Reactors []reactorLifecycleView `json:"reactors"`
}

type reactorLifecycleView struct {
	Reactor   reactorView   `json:"reactor"`
	Lifecycle lifecycleView `json:"lifecycle"`
}

type healthView battlefield.Ratio

func (v healthView) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]int{
		"current": v.Current,
		"maximum": v.Maximum,
	})
}

type profile struct {
	Warrior      warriorView            `json:"warrior"`
	Health       healthView             `json:"health"`
	Damage       int                    `json:"damage"`
	Defense      int                    `json:"defense"`
	CriticalOdds int                    `json:"critical_odds"`
	CriticalLoss int                    `json:"critical_loss"`
	Speed        int                    `json:"speed"`
	Reactors     []reactorLifecycleView `json:"reactors"`
}

func newProfile(warrior battlefield.Warrior, ec battlefield.EvaluationContext) *profile {
	return &profile{
		Warrior:      newWarriorView(warrior),
		Health:       healthView(warrior.Health()),
		Damage:       warrior.Component(battlefield.Damage, ec),
		Defense:      warrior.Component(battlefield.Defense, ec),
		CriticalOdds: warrior.Component(battlefield.CriticalOdds, ec),
		CriticalLoss: warrior.Component(battlefield.CriticalLoss, ec),
		Speed:        warrior.Component(battlefield.Speed, ec),
		Reactors: functional.MapSlice(func(reactor battlefield.Reactor) reactorLifecycleView {
			return reactorLifecycleView{
				newReactorView(reactor),
				lifecycleView(reactor.(*battlefield.FatReactor).Lifecycle()),
			}
		}, warrior.Buffs()),
	}
}

func createProfiles(ec battlefield.EvaluationContext) []*profile {
	ps := make([]*profile, 0, len(ec.Warriors()))
	for _, w := range ec.Warriors() {
		if w.Health().Current > 0 {
			ps = append(ps, newProfile(w, ec))
		}
	}

	return ps
}
