package controller

import (
	"encoding/json"

	"github.com/farseeingnorthwest/battleground.go/functional"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/gofiber/fiber/v2"
)

type SkillController struct {
	repo SkillRepository
}

type SkillRepository interface {
	Find(ids ...int) ([]storage.SkillMeta, error)
	Create(skill *storage.Skill) error
	Get(id int) (*storage.Skill, error)
	Update(skill *storage.Skill) error
	Delete(id int, force bool) error
}

func NewSkillController(repo SkillRepository) SkillController {
	return SkillController{repo}
}

func (c SkillController) Mount(app *fiber.App) {
	app.Get("/skills", c.GetSkills)
	app.Post("/skills", c.CreateSkill)
	app.Get("/skills/:id", c.GetSkill)
	app.Put("/skills/:id", c.UpdateSkill)
	app.Delete("/skills/:id", c.DeleteSkill)
}

func (c SkillController) GetSkills(fc *fiber.Ctx) error {
	skills, err := c.repo.Find()
	if err != nil {
		return err
	}

	return fc.JSON(functional.MapSlice(newSkillMetaView, skills))
}

func (c SkillController) CreateSkill(fc *fiber.Ctx) error {
	var form struct {
		Name    string
		Reactor battlefield.FatReactorFile
	}
	if err := fc.BodyParser(&form); err != nil {
		return err
	}

	skill := storage.Skill{
		SkillMeta: storage.SkillMeta{
			Name: form.Name,
		},
		Reactor: (*storage.Reactor)(form.Reactor.FatReactor),
	}
	if err := c.repo.Create(&skill); err != nil {
		return err
	}

	return fc.JSON(skillView(skill))
}

func (c SkillController) GetSkill(fc *fiber.Ctx) error {
	id, err := fc.ParamsInt("id")
	if err != nil {
		return err
	}
	skill, err := c.repo.Get(id)
	if err != nil {
		return err
	}

	return fc.JSON((*skillView)(skill))
}

func (c SkillController) UpdateSkill(fc *fiber.Ctx) error {
	id, err := fc.ParamsInt("id")
	if err != nil {
		return err
	}

	var form struct {
		Name    string
		Reactor battlefield.FatReactorFile
	}
	if err := fc.BodyParser(&form); err != nil {
		return err
	}

	skill := storage.Skill{
		SkillMeta: storage.SkillMeta{
			ID:   id,
			Name: form.Name,
		},
		Reactor: (*storage.Reactor)(form.Reactor.FatReactor),
	}
	if err := c.repo.Update(&skill); err != nil {
		return err
	}

	return fc.JSON(skillView(skill))
}

func (c SkillController) DeleteSkill(fc *fiber.Ctx) error {
	id, err := fc.ParamsInt("id")
	if err != nil {
		return err
	}
	if err := c.repo.Delete(id, false); err != nil {
		return err
	}

	return fc.SendStatus(fiber.StatusNoContent)
}

type skillMetaView storage.SkillMeta

func newSkillMetaView(skill storage.SkillMeta) skillMetaView {
	return skillMetaView(skill)
}

func (s skillMetaView) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":   s.ID,
		"name": s.Name,
	})
}

type skillView storage.Skill

func newSkillView(skill storage.Skill) skillView {
	return skillView(skill)
}

func (v skillView) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"id":      v.ID,
		"name":    v.Name,
		"reactor": (*battlefield.FatReactor)(v.Reactor),
	})
}

type bySkillID []storage.SkillMeta

func (s bySkillID) Len() int                           { return len(s) }
func (s bySkillID) Get(i int) (int, storage.SkillMeta) { return s[i].ID, s[i] }
