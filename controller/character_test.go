package controller_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/farseeingnorthwest/battleground.go/controller"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCharacterController_GetCharacters(t *testing.T) {
	r := new(mockCharacterRepository)
	sr := new(mockSkillRepository)
	r.On("Find", []int(nil)).Return([]storage.Character{
		{
			ID:           1,
			Name:         "Oda",
			Damage:       10,
			Defense:      5,
			CriticalOdds: 10,
			CriticalLoss: 200,
			Health:       100,
			Speed:        10,
			Skills: map[int]storage.SkillMeta{
				1: {
					ID:   1,
					Name: "Normal Attack",
				},
			},
		},
	}, nil)

	app := fiber.New()
	NewCharacterController(r, sr).Mount(app)
	req := httptest.NewRequest("GET", "/characters", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Oda")
	assert.Contains(t, string(body), "Normal Attack")
}

func TestCharacterController_CreateCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
	sr := new(mockSkillRepository)
	r.On("Create", &storage.Character{
		Name:         "Toy",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
		Skills: map[int]storage.SkillMeta{
			1: {
				ID:   2,
				Name: "Sleep",
			},
		},
	}).Run(func(args mock.Arguments) {
		args.Get(0).(*storage.Character).ID = 1
	}).Return(nil)
	sr.On("Find", []int{2}).Return([]storage.SkillMeta{
		{
			ID:   2,
			Name: "Sleep",
		},
	}, nil)

	app := fiber.New()
	NewCharacterController(r, sr).Mount(app)
	req := httptest.NewRequest("POST", "/characters", strings.NewReader(
		`{"name":"Toy","damage":9,"defense":4,"critical_odds":10,"critical_loss":150,"health":80,"speed":9,"skills":{"1":2}}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	sr.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Toy")
	assert.Contains(t, string(body), "Sleep")
}

func TestCharacterController_GetCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
	sr := new(mockSkillRepository)
	r.On("Get", 1).Return(&storage.Character{
		ID:           1,
		Name:         "Oda",
		Damage:       10,
		Defense:      5,
		CriticalOdds: 10,
		CriticalLoss: 200,
		Health:       100,
		Speed:        10,
	}, nil)

	app := fiber.New()
	NewCharacterController(r, sr).Mount(app)
	req := httptest.NewRequest("GET", "/characters/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Oda")
}

func TestCharacterController_UpdateCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
	sr := new(mockSkillRepository)
	r.On("Update", &storage.Character{
		ID:           1,
		Name:         "Oda",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
		Skills: map[int]storage.SkillMeta{
			4: {
				ID:   2,
				Name: "Sleep",
			},
		},
	}).Return(nil)
	sr.On("Find", []int{2}).Return([]storage.SkillMeta{
		{
			ID:   2,
			Name: "Sleep",
		},
	}, nil)

	app := fiber.New()
	NewCharacterController(r, sr).Mount(app)
	req := httptest.NewRequest("PUT", "/characters/1", strings.NewReader(
		`{"name":"Oda","damage":9,"defense":4,"critical_odds":10,"critical_loss":150,"health":80,"speed":9,"skills":{"4":2}}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	sr.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "Oda")
	assert.Contains(t, string(body), "Sleep")
}

func TestCharacterController_DeleteCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
	sr := new(mockSkillRepository)
	r.On("Delete", 1).Return(nil)

	app := fiber.New()
	NewCharacterController(r, sr).Mount(app)
	req := httptest.NewRequest("DELETE", "/characters/1", nil)
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

type mockCharacterRepository struct {
	mock.Mock
}

func (r *mockCharacterRepository) Find(ids ...int) ([]storage.Character, error) {
	args := r.Called(ids)
	return args.Get(0).([]storage.Character), args.Error(1)
}

func (r *mockCharacterRepository) Create(character *storage.Character) error {
	args := r.Called(character)
	return args.Error(0)
}

func (r *mockCharacterRepository) Get(id int) (*storage.Character, error) {
	args := r.Called(id)
	return args.Get(0).(*storage.Character), args.Error(1)
}

func (r *mockCharacterRepository) Update(character *storage.Character) error {
	args := r.Called(character)
	return args.Error(0)
}

func (r *mockCharacterRepository) Delete(id int, force bool) error {
	args := r.Called(id)
	return args.Error(0)
}
