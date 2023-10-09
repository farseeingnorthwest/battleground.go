package controller_test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/farseeingnorthwest/battleground.go/controller"
	"github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/farseeingnorthwest/playground/battlefield/v2/examples"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestBattleController_CreateBattle(t *testing.T) {
	sch, err := jsonschema.Compile("battle.schema.json")
	assert.NoError(t, err)

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
		{
			SkillMeta: storage.SkillMeta{
				ID:   3,
				Name: "#1-1",
			},
			Reactor: (*storage.Reactor)(examples.Special[0][0]),
		},
		{
			SkillMeta: storage.SkillMeta{
				ID:   4,
				Name: "#1-2",
			},
			Reactor: (*storage.Reactor)(examples.Special[0][1]),
		},
		{
			SkillMeta: storage.SkillMeta{
				ID:   5,
				Name: "#1-3",
			},
			Reactor: (*storage.Reactor)(examples.Special[0][2]),
		},
		{
			SkillMeta: storage.SkillMeta{
				ID:   6,
				Name: "#1-4",
			},
			Reactor: (*storage.Reactor)(examples.Special[0][3]),
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
			Health:       200,
			Speed:        10,
			Skills: map[int]storage.SkillMeta{
				0: {
					ID:   1,
					Name: "Normal Attack",
				},
				1: {
					ID:   3,
					Name: "#1-1",
				},
				2: {
					ID:   4,
					Name: "#1-2",
				},
				3: {
					ID:   5,
					Name: "#1-3",
				},
				4: {
					ID:   6,
					Name: "#1-4",
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
			Health:       180,
			Speed:        11,
			Skills: map[int]storage.SkillMeta{
				0: {
					ID:   1,
					Name: "Normal Attack",
				},
			},
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

	var v any
	assert.NoError(t, json.Unmarshal(body, &v))
	assert.NoError(t, sch.Validate(v))
}
