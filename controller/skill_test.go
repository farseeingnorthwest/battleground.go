package controller_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/farseeingnorthwest/battleground.go/controller"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/farseeingnorthwest/playground/battlefield/v2/examples"
	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSkillController_GetSkills(t *testing.T) {
	r := new(mockSkillRepository)
	r.On("Find", []int(nil)).Return([]storage.SkillMeta{
		{
			ID:   1,
			Name: "Normal Attack",
		},
	}, nil)

	app := fiber.New()
	NewSkillController(r).Mount(app)
	req := httptest.NewRequest("GET", "/skills", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `[{"id":1,"name":"Normal Attack"}]`, string(body))
}

func TestSkillController_CreateSkill(t *testing.T) {
	r := new(mockSkillRepository)
	r.On("Create", &storage.Skill{
		SkillMeta: storage.SkillMeta{
			Name: "Sleep",
		},
		Reactor: (*storage.Reactor)(examples.Effect["Sleep"]),
	}).Run(func(args mock.Arguments) {
		args.Get(0).(*storage.Skill).ID = 1
	}).Return(nil)

	app := fiber.New()
	NewSkillController(r).Mount(app)
	req := httptest.NewRequest("POST", "/skills", strings.NewReader(
		`{"name":"Sleep","reactor":{"tags":[{"_kind":"exclusion_group","index":0},{"_kind":"priority","index":10},{"_kind":"label","text":"Sleep"}],"capacity":{"count":1,"when":[{"signal":"round_end"},{"if":[{"_kind":"verb","verb":"attack"},{"_kind":"current_is_target"}],"signal":"post_action"}]},"cases":[{"when":{"signal":"launch"},"then":{"_kind":"sequence","do":[]}}]}}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(functional.First(io.ReadAll(resp.Body))), "Sleep")
}

func TestSkillController_GetSkill(t *testing.T) {
	r := new(mockSkillRepository)
	r.On("Get", 1).Return(&storage.Skill{
		SkillMeta: storage.SkillMeta{
			ID:   1,
			Name: "Normal Attack",
		},
		Reactor: (*storage.Reactor)(examples.Regular[0]),
	}, nil)

	app := fiber.New()
	NewSkillController(r).Mount(app)
	req := httptest.NewRequest("GET", "/skills/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(functional.First(io.ReadAll(resp.Body))), "NormalAttack")
}

func TestSkillController_UpdateSkill(t *testing.T) {
	r := new(mockSkillRepository)
	r.On("Update", &storage.Skill{
		SkillMeta: storage.SkillMeta{
			ID:   1,
			Name: "Sleep",
		},
		Reactor: (*storage.Reactor)(examples.Effect["Sleep"]),
	}).Return(nil)

	app := fiber.New()
	NewSkillController(r).Mount(app)
	req := httptest.NewRequest("PUT", "/skills/1", strings.NewReader(
		`{"name":"Sleep","reactor":{"tags":[{"_kind":"exclusion_group","index":0},{"_kind":"priority","index":10},{"_kind":"label","text":"Sleep"}],"capacity":{"count":1,"when":[{"signal":"round_end"},{"if":[{"_kind":"verb","verb":"attack"},{"_kind":"current_is_target"}],"signal":"post_action"}]},"cases":[{"when":{"signal":"launch"},"then":{"_kind":"sequence","do":[]}}]}}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(functional.First(io.ReadAll(resp.Body))), "Sleep")
}

func TestSkillController_DeleteSkill(t *testing.T) {
	r := new(mockSkillRepository)
	r.On("Delete", 1).Return(nil)

	app := fiber.New()
	NewSkillController(r).Mount(app)
	req := httptest.NewRequest("DELETE", "/skills/1", nil)
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

type mockSkillRepository struct {
	mock.Mock
}

func (r *mockSkillRepository) Find(ids ...int) ([]storage.SkillMeta, error) {
	args := r.Called(ids)
	return args.Get(0).([]storage.SkillMeta), args.Error(1)
}

func (r *mockSkillRepository) Create(skill *storage.Skill) error {
	args := r.Called(skill)
	return args.Error(0)
}

func (r *mockSkillRepository) Get(id int) (*storage.Skill, error) {
	args := r.Called(id)
	return args.Get(0).(*storage.Skill), args.Error(1)
}

func (r *mockSkillRepository) Update(skill *storage.Skill) error {
	args := r.Called(skill)
	return args.Error(0)
}

func (r *mockSkillRepository) Delete(id int, force bool) error {
	args := r.Called(id)
	return args.Error(0)
}
