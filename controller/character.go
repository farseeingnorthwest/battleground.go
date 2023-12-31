package controller

import (
	"encoding/json"

	"github.com/farseeingnorthwest/battleground.go/functional"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/gofiber/fiber/v2"
)

type CharacterController struct {
	repo      CharacterRepository
	skillRepo SkillRepository
}

type CharacterRepository interface {
	Find(...int) ([]storage.Character, error)
	Create(*storage.Character) error
	Get(int) (*storage.Character, error)
	Update(*storage.Character) error
	Delete(int, bool) error
}

func NewCharacterController(repo CharacterRepository, skillRepo SkillRepository) CharacterController {
	return CharacterController{repo, skillRepo}
}

func (c CharacterController) Mount(router fiber.Router) {
	router.Get("/characters", c.GetCharacters)
	router.Post("/characters", c.CreateCharacter)
	router.Get("/characters/:id", c.GetCharacter)
	router.Put("/characters/:id", c.UpdateCharacter)
	router.Delete("/characters/:id", c.DeleteCharacter)
}

func (c CharacterController) GetCharacters(fc *fiber.Ctx) error {
	characters, err := c.repo.Find()
	if err != nil {
		return err
	}

	return fc.JSON(functional.MapSlice(newCharacterView, characters))
}

func (c CharacterController) CreateCharacter(fc *fiber.Ctx) error {
	character, err := c.form(fc, false)
	if err != nil {
		return err
	}
	if err := c.repo.Create(character); err != nil {
		return err
	}

	return fc.JSON((*characterView)(character))
}

func (c CharacterController) GetCharacter(fc *fiber.Ctx) error {
	id, err := fc.ParamsInt("id")
	if err != nil {
		return err
	}

	character, err := c.repo.Get(id)
	if err != nil {
		return err
	}

	return fc.JSON((*characterView)(character))
}

func (c CharacterController) UpdateCharacter(fc *fiber.Ctx) error {
	character, err := c.form(fc, true)
	if err != nil {
		return err
	}
	if err = c.repo.Update(character); err != nil {
		return err
	}

	return fc.JSON((*characterView)(character))
}

func (c CharacterController) DeleteCharacter(fc *fiber.Ctx) error {
	id, err := fc.ParamsInt("id")
	if err != nil {
		return err
	}

	if err := c.repo.Delete(id, false); err != nil {
		return err
	}

	return fc.SendStatus(fiber.StatusNoContent)
}

func (c CharacterController) form(fc *fiber.Ctx, withID bool) (*storage.Character, error) {
	var id int
	if withID {
		var err error
		if id, err = fc.ParamsInt("id"); err != nil {
			return nil, err
		}
	}

	form := struct {
		Name         string
		Damage       int
		Defense      int
		CriticalOdds int `json:"critical_odds"`
		CriticalLoss int `json:"critical_loss"`
		Health       int
		Speed        int
		Skills       map[int]int
	}{
		Skills: make(map[int]int),
	}
	if err := fc.BodyParser(&form); err != nil {
		return nil, err
	}
	var skills map[int]storage.SkillMeta
	if len(form.Skills) > 0 {
		skillMetas, err := c.skillRepo.Find(functional.Values(form.Skills)...)
		if err != nil {
			return nil, err
		}

		metas := functional.Tabulate[int, storage.SkillMeta](bySkillMetaID(skillMetas))
		skills = make(map[int]storage.SkillMeta)
		for slot, id := range form.Skills {
			skills[slot] = metas[id]
		}
	}

	return &storage.Character{
		ID:           id,
		Name:         form.Name,
		Damage:       form.Damage,
		Defense:      form.Defense,
		CriticalOdds: form.CriticalOdds,
		CriticalLoss: form.CriticalLoss,
		Health:       form.Health,
		Speed:        form.Speed,
		Skills:       skills,
	}, nil
}

type characterView storage.Character

func newCharacterView(character storage.Character) characterView {
	return characterView(character)
}

func (c characterView) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":            c.ID,
		"name":          c.Name,
		"damage":        c.Damage,
		"defense":       c.Defense,
		"critical_odds": c.CriticalOdds,
		"critical_loss": c.CriticalLoss,
		"health":        c.Health,
		"speed":         c.Speed,
		"skills":        functional.MapValues(newSkillMetaView, c.Skills),
	})
}

type byCharacterID []storage.Character

func (s byCharacterID) Len() int                           { return len(s) }
func (s byCharacterID) Get(i int) (int, storage.Character) { return s[i].ID, s[i] }
