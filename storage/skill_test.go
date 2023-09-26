package storage_test

import (
	"testing"

	. "github.com/farseeingnorthwest/battleground.go/storage"
	b "github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/stretchr/testify/assert"
)

func TestSkillRepository_Find(t *testing.T) {
	for _, tt := range []struct {
		ids    []int
		skills []SkillMeta
	}{
		{nil, []SkillMeta{{1, "Normal Attack"}, {2, "Sleep"}}},
		{[]int{1}, []SkillMeta{{1, "Normal Attack"}}},
	} {
		t.Run("", func(t *testing.T) {
			loadFixtures(t)

			r := NewSkillRepository(db)
			skills, err := r.Find(tt.ids...)

			assert.NoError(t, err)
			assert.Equal(t, tt.skills, skills)
		})
	}
}

func TestSkillRepository_FindEx(t *testing.T) {
	loadFixtures(t)

	r := NewSkillRepository(db)
	skills, err := r.FindEx()

	assert.NoError(t, err)
	assert.Len(t, skills, 2)
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
		SkillMeta: SkillMeta{
			Name: "Taunt",
		},
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
		SkillMeta: SkillMeta{
			ID:   1,
			Name: "Taunt",
		},
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
		t.Run("", func(t *testing.T) {
			loadFixtures(t)

			r := NewSkillRepository(db)
			err := r.Delete(tt.id, tt.force)
			if tt.ok {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}

			skills, err := r.Find()
			assert.NoError(t, err)
			assert.Len(t, skills, tt.count)
		})
	}
}
