package controller_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/farseeingnorthwest/battleground.go/controller"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/farseeingnorthwest/playground/battlefield/v2/examples"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestBattleController_CreateBattle(t *testing.T) {
	r := new(mockCharacterRepository)
	sr := new(mockSkillRepository)
	sr.On("FindEx", []int(nil)).Return([]storage.Skill{
		{
			SkillMeta: storage.SkillMeta{
				ID:   1,
				Name: "Normal Attack",
			},
			Reactor: (*storage.Reactor)(examples.Regular[0]),
		},
		{
			SkillMeta: storage.SkillMeta{
				ID:   2,
				Name: "Element Theory",
			},
			Reactor: (*storage.Reactor)(examples.Regular[3]),
		},
	}, nil)
	r.On("Find", []int{1}).Return([]storage.Character{
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
				0: {
					ID:   1,
					Name: "Normal Attack",
				},
			},
		},
	}, nil)
	r.On("Find", []int{2}).Return([]storage.Character{
		{
			ID:           2,
			Name:         "Ueno",
			Damage:       9,
			Defense:      4,
			CriticalOdds: 20,
			CriticalLoss: 200,
			Health:       90,
			Speed:        11,
		},
	}, nil)

	app := fiber.New()
	controller.NewBattleController(r, sr).Mount(app)
	req := httptest.NewRequest("POST", "/battles", strings.NewReader(
		`{"left":{"0":1},"right":{"0":2},"ground":[2]}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	r.AssertExpectations(t)
	sr.AssertExpectations(t)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"winner":"Left"}`, string(body))
}
