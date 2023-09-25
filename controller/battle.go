package controller

import (
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

func (c BattleController) Mount(app *fiber.App) {
	app.Post("/battles", c.CreateBattle)
}

func (c BattleController) CreateBattle(fc *fiber.Ctx) error {
	form := struct {
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
	ob := observer{battlefield.NewTagSet(battlefield.Priority(1000000)), nil}
	f := battlefield.NewBattleField(
		append(left, right...),
		append(
			functional.MapSlice(
				func(id int) battlefield.Reactor {
					return skills[id].Spawn()
				},
				form.Ground,
			),
			&ob,
		)...,
	)
	f.Run()

	return fc.JSON(fiber.Map{
		"winner": ob.winner.Side().String(),
	})
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
			},
			side,
			p,
			battlefield.WarriorSkills(functional.MapSlice(
				func(meta storage.SkillMeta) battlefield.Reactor {
					return skills[meta.ID].Spawn()
				},
				functional.Values(chars[id].Skills),
			)...),
		))
	}

	return warriors, nil
}

type observer struct {
	battlefield.TagSet
	winner battlefield.Warrior
}

func (o *observer) React(signal battlefield.Signal, ec battlefield.EvaluationContext) {
	switch signal := signal.(type) {
	case *battlefield.PostActionSignal:
		_, source, _ := signal.Action().Script().Source()
		o.winner = source.(battlefield.Warrior)
	}
}

func (o *observer) Active() bool {
	return true
}
