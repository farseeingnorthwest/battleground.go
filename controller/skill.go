package controller

import (
	"encoding/json"

	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
	"github.com/gofiber/fiber/v2"
)

type SkillController struct {
	repo SkillRepository
}

type SkillRepository interface {
	Find(query *storage.SkillQuery) ([]storage.Skill, error)
	Create(skill *storage.Skill) error
	Get(id int) (*storage.Skill, error)
	Update(skill *storage.Skill) error
	Delete(id int) error
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
	skills, err := c.repo.Find(&storage.SkillQuery{})
	if err != nil {
		return err
	}

	return fc.JSON(functional.Map(newSkillView)(skills))
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
		Name:    form.Name,
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
		ID:      id,
		Name:    form.Name,
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
	if err := c.repo.Delete(id); err != nil {
		return err
	}

	return fc.SendStatus(fiber.StatusNoContent)
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
