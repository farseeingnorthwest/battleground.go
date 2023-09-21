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
	assert.Len(t, characters, 1)
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
	}, character)
}

func TestCharacterRepository_Create(t *testing.T) {
	loadFixtures(t)

	r := NewCharacterRepository(db)
	toyotomi := Character{
		Name:         "Toyotomi",
		Damage:       9,
		Defense:      4,
		CriticalOdds: 10,
		CriticalLoss: 150,
		Health:       80,
		Speed:        9,
	}
	err := r.Create(&toyotomi)
	assert.NoError(t, err)
	assert.NotEmpty(t, toyotomi.ID)

	character, err := r.Get(toyotomi.ID)
	assert.NoError(t, err)
	assert.Equal(t, &toyotomi, character)
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
	}, character)
}

func TestCharacterRepository_Delete(t *testing.T) {
	loadFixtures(t)

	r := NewCharacterRepository(db)
	err := r.Delete(1)
	assert.NoError(t, err)

	characters, err := r.Find()
	assert.NoError(t, err)
	assert.Len(t, characters, 0)
}
