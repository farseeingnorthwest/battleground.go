package controller_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/farseeingnorthwest/battleground.go/controller"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/farseeingnorthwest/playground/battlefield/v2/functional"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCharacterController_GetCharacters(t *testing.T) {
	r := new(mockCharacterRepository)
	r.On("Find").Return([]storage.Character{
		{
			ID:           1,
			Name:         "Oda",
			Damage:       10,
			Defense:      5,
			CriticalOdds: 10,
			CriticalLoss: 200,
			Health:       100,
			Speed:        10,
		},
	}, nil)

	app := fiber.New()
	NewCharacterController(r).Mount(app)
	req := httptest.NewRequest("GET", "/characters", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(functional.First(io.ReadAll(resp.Body))), "Oda")
}

func TestCharacterController_CreateCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
	r.On("Create", &storage.Character{
		Name:         "Toyotomi",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
	}).Run(func(args mock.Arguments) {
		args.Get(0).(*storage.Character).ID = 1
	}).Return(nil)

	app := fiber.New()
	NewCharacterController(r).Mount(app)
	req := httptest.NewRequest("POST", "/characters", strings.NewReader(
		`{"name":"Toyotomi","damage":9,"defense":4,"critical_odds":10,"critical_loss":150,"health":80,"speed":9}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(functional.First(io.ReadAll(resp.Body))), "Toyotomi")
}

func TestCharacterController_GetCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
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
	NewCharacterController(r).Mount(app)
	req := httptest.NewRequest("GET", "/characters/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(functional.First(io.ReadAll(resp.Body))), "Oda")
}

func TestCharacterController_UpdateCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
	r.On("Update", &storage.Character{
		ID:           1,
		Name:         "Oda",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
	}).Return(nil)

	app := fiber.New()
	NewCharacterController(r).Mount(app)
	req := httptest.NewRequest("PUT", "/characters/1", strings.NewReader(
		`{"name":"Oda","damage":9,"defense":4,"critical_odds":10,"critical_loss":150,"health":80,"speed":9}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Contains(t, string(functional.First(io.ReadAll(resp.Body))), "Oda")
}

func TestCharacterController_DeleteCharacter(t *testing.T) {
	r := new(mockCharacterRepository)
	r.On("Delete", 1).Return(nil)

	app := fiber.New()
	NewCharacterController(r).Mount(app)
	req := httptest.NewRequest("DELETE", "/characters/1", nil)
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
}

type mockCharacterRepository struct {
	mock.Mock
}

func (r *mockCharacterRepository) Find() ([]storage.Character, error) {
	args := r.Called()
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

func (r *mockCharacterRepository) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}
