package storage_test

import (
	"testing"

	. "github.com/farseeingnorthwest/battleground.go/storage"
	"github.com/stretchr/testify/assert"
)

func TestCharacterRepository_Find(t *testing.T) {
	loadFixtures(t)

	r := NewCharacterRepository(db)
	characters, err := r.Find()

	assert.NoError(t, err)
	assert.Equal(t, []Character{
		{
			ID:           1,
			Name:         "Oda",
			Damage:       10,
			Defense:      5,
			CriticalOdds: 10,
			CriticalLoss: 200,
			Health:       100,
			Speed:        10,
			Skills: map[int]SkillMeta{
				1: {ID: 1, Name: "Normal Attack"},
			},
		},
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
	}, characters)
}

func TestCharacterRepository_Get(t *testing.T) {
	loadFixtures(t)

	r := NewCharacterRepository(db)
	character, err := r.Get(1)

	assert.NoError(t, err)
	assert.Equal(t, &Character{
		ID:           1,
		Name:         "Oda",
		Damage:       10,
		Defense:      5,
		CriticalOdds: 10,
		CriticalLoss: 200,
		Health:       100,
		Speed:        10,
		Skills: map[int]SkillMeta{
			1: {ID: 1, Name: "Normal Attack"},
		},
	}, character)
}

func TestCharacterRepository_Create(t *testing.T) {
	loadFixtures(t)

	r := NewCharacterRepository(db)
	toy := Character{
		Name:         "Toy",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
		Skills: map[int]SkillMeta{
			1: {ID: 1, Name: "Normal Attack"},
		},
	}
	err := r.Create(&toy)
	assert.NoError(t, err)
	assert.NotEmpty(t, toy.ID)

	character, err := r.Get(toy.ID)
	assert.NoError(t, err)
	assert.Equal(t, &toy, character)
}

func TestCharacterRepository_Update(t *testing.T) {
	loadFixtures(t)

	r := NewCharacterRepository(db)
	err := r.Update(&Character{
		ID:           1,
		Name:         "Oda",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
		Skills: map[int]SkillMeta{
			4: {ID: 2, Name: "Sleep"},
		},
	})
	assert.NoError(t, err)

	character, err := r.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, &Character{
		ID:           1,
		Name:         "Oda",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
		Skills: map[int]SkillMeta{
			4: {ID: 2, Name: "Sleep"},
		},
	}, character)
}

func TestCharacterRepository_Delete(t *testing.T) {
	for _, tt := range []struct {
		id    int
		force bool
		ok    bool
		count int
	}{
		{1, false, false, 2},
		{1, true, true, 1},
		{2, false, true, 1},
	} {
		loadFixtures(t)

		r := NewCharacterRepository(db)
		err := r.Delete(tt.id, tt.force)
		if tt.ok {
			assert.NoError(t, err)
		} else {
			assert.Error(t, err)
		}

		characters, err := r.Find()
		assert.NoError(t, err)
		assert.Len(t, characters, tt.count)
	}
}
