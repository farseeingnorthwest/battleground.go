package storage_test

import (
	"testing"

	. "github.com/farseeingnorthwest/battleground.go/storage"
	b "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestSkillRepository_List(t *testing.T) {
	loadFixtures(t)

	r := NewSkillRepository(db)
	skills, err := r.Find(&SkillQuery{})

	assert.NoError(t, err)
	assert.Len(t, skills, 1)
}

func TestSkillRepository_Get(t *testing.T) {
	loadFixtures(t)

	r := NewSkillRepository(db)
	skill, err := r.Get(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, skill.ID)
	assert.Equal(t, "Normal Attack", skill.Name)
	assert.Contains(t, skill.Reactor.Tags(), b.Label("NormalAttack"))
}

func TestSkillRepository_Create(t *testing.T) {
	loadFixtures(t)

	r := NewSkillRepository(db)
	taunt := Skill{
		Name: "Taunt",
		Reactor: (*Reactor)(b.NewFatReactor(
			b.FatTags(b.Label("Taunt")),
			b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 2),
		)),
	}
	err := r.Create(&taunt)
	assert.NoError(t, err)
	assert.NotEmpty(t, taunt.ID)

	skill, err := r.Get(taunt.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Taunt", skill.Name)
	assert.Contains(t, skill.Reactor.Tags(), b.Label("Taunt"))
}

func TestSkillRepository_Update(t *testing.T) {
	loadFixtures(t)

	r := NewSkillRepository(db)
	err := r.Update(&Skill{
		ID:   1,
		Name: "Taunt",
		Reactor: (*Reactor)(b.NewFatReactor(
			b.FatTags(b.Label("Taunt")),
			b.FatCapacity(b.NewSignalTrigger(&b.RoundEndSignal{}), 2),
		)),
	})
	assert.NoError(t, err)

	skill, err := r.Get(1)
	assert.NoError(t, err)
	assert.Equal(t, "Taunt", skill.Name)
	assert.Contains(t, skill.Reactor.Tags(), b.Label("Taunt"))
}

func TestSkillRepository_Delete(t *testing.T) {
	loadFixtures(t)

	r := NewSkillRepository(db)
	err := r.Delete(1)
	assert.NoError(t, err)

	skills, err := r.Find(&SkillQuery{})
	assert.NoError(t, err)
	assert.Len(t, skills, 0)
}
