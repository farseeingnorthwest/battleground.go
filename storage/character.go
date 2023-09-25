package storage

import "github.com/jmoiron/sqlx"

type Character struct {
	ID           int
	Name         string
	Damage       int
	Defense      int
	CriticalOdds int `db:"critical_odds"`
	CriticalLoss int `db:"critical_loss"`
	Health       int
	Speed        int
	Skills       map[int]SkillMeta
}

type CharacterSkill struct {
	CharacterID int `db:"character_id"`
	Slot        int
	SkillMeta
}

type CharacterRepository struct {
	db *sqlx.DB
}

func NewCharacterRepository(db *sqlx.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (r CharacterRepository) Find(ids ...int) ([]Character, error) {
	var characters []Character
	if len(ids) == 0 {
		if err := r.db.Select(&characters, "SELECT * FROM characters ORDER BY id"); err != nil {
			return nil, err
		}
	} else {
		query, args, err := sqlx.In("SELECT * FROM characters WHERE id IN (?) ORDER BY id", ids)
		if err != nil {
			return nil, err
		}

		if err := r.db.Select(&characters, query, args...); err != nil {
			return nil, err
		}
	}

	return r.getAllCharacterSkills(characters)
}

func (r CharacterRepository) Get(id int) (*Character, error) {
	var character Character
	if err := r.db.Get(&character, "SELECT * FROM characters WHERE id = $1", id); err != nil {
		return nil, err
	}

	return r.getCharacterSkills(&character)
}

func (r CharacterRepository) Create(character *Character) error {
	if err := r.db.Get(
		character, `
INSERT INTO
    characters (name, damage, defense, critical_odds, critical_loss, health, speed)
VALUES
    ($1, $2, $3, $4, $5, $6, $7)
RETURNING
    *
`,
		character.Name,
		character.Damage,
		character.Defense,
		character.CriticalOdds,
		character.CriticalLoss,
		character.Health,
		character.Speed,
	); err != nil {
		return err
	}

	return r.saveCharacterSkills(character)
}

func (r CharacterRepository) Update(character *Character) error {
	if err := r.db.Get(
		character, `
UPDATE
    characters
SET
    name = $1,
    damage = $2,
    defense = $3,
    critical_odds = $4,
    critical_loss = $5,
    health = $6,
    speed = $7
WHERE
    id = $8
RETURNING *
`,
		character.Name,
		character.Damage,
		character.Defense,
		character.CriticalOdds,
		character.CriticalLoss,
		character.Health,
		character.Speed,
		character.ID,
	); err != nil {
		return err
	}

	return r.saveCharacterSkills(character)
}

func (r CharacterRepository) Delete(id int, force bool) error {
	if force {
		if err := r.removeCharacterSkills(id); err != nil {
			return err
		}
	}
	if _, err := r.db.Exec("DELETE FROM characters WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}

func (r CharacterRepository) getAllCharacterSkills(characters []Character) ([]Character, error) {
	var skills []CharacterSkill
	if err := r.db.Select(&skills, `
SELECT
    character_id, slot, id, name
FROM
    character_skills c JOIN
        skills s ON c.skill_id = s.id
ORDER BY
    character_id, slot
`,
	); err != nil {
		return nil, err
	}

	var byCharacter = make(map[int]map[int]SkillMeta)
	for _, skill := range skills {
		if _, ok := byCharacter[skill.CharacterID]; !ok {
			byCharacter[skill.CharacterID] = make(map[int]SkillMeta)
		}
		byCharacter[skill.CharacterID][skill.Slot] = skill.SkillMeta
	}
	for i := range characters {
		characters[i].Skills = byCharacter[characters[i].ID]
	}

	return characters, nil
}

func (r CharacterRepository) getCharacterSkills(character *Character) (*Character, error) {
	var skills []CharacterSkill
	if err := r.db.Select(&skills, `
SELECT
    character_id, slot, id, name
FROM
    character_skills c JOIN
        skills s ON c.skill_id = s.id
WHERE
    character_id = $1
ORDER BY
    slot
`,
		character.ID); err != nil {
		return nil, err
	}

	character.Skills = make(map[int]SkillMeta)
	for _, skill := range skills {
		character.Skills[skill.Slot] = skill.SkillMeta
	}

	return character, nil
}

func (r CharacterRepository) saveCharacterSkills(character *Character) error {
	if err := r.removeCharacterSkills(character.ID); err != nil {
		return err
	}

	for slot, skill := range character.Skills {
		if _, err := r.db.Exec(
			"INSERT INTO character_skills (character_id, slot, skill_id) VALUES ($1, $2, $3)",
			character.ID, slot, skill.ID); err != nil {
			return err
		}
	}

	return nil
}

func (r CharacterRepository) removeCharacterSkills(id int) error {
	if _, err := r.db.Exec("DELETE FROM character_skills WHERE character_id = $1", id); err != nil {
		return err
	}

	return nil
}
