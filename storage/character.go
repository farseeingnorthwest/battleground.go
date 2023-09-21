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
}

type CharacterRepository struct {
	db *sqlx.DB
}

func NewCharacterRepository(db *sqlx.DB) *CharacterRepository {
	return &CharacterRepository{db: db}
}

func (r CharacterRepository) Find() (characters []Character, err error) {
	err = r.db.Select(&characters, "SELECT * FROM characters ORDER BY id")
	return
}

func (r CharacterRepository) Get(id int) (*Character, error) {
	var character Character
	if err := r.db.Get(&character, "SELECT * FROM characters WHERE id = $1", id); err != nil {
		return nil, err
	}

	return &character, nil
}

func (r CharacterRepository) Create(character *Character) error {
	return r.db.Get(character, "INSERT INTO characters (name, damage, defense, critical_odds, critical_loss, health, speed) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *",
		character.Name, character.Damage, character.Defense, character.CriticalOdds, character.CriticalLoss, character.Health, character.Speed)
}

func (r CharacterRepository) Update(character *Character) error {
	return r.db.Get(character, "UPDATE characters SET name = $1, damage = $2, defense = $3, critical_odds = $4, critical_loss = $5, health = $6, speed = $7 WHERE id = $8 RETURNING *",
		character.Name, character.Damage, character.Defense, character.CriticalOdds, character.CriticalLoss, character.Health, character.Speed, character.ID)
}

func (r CharacterRepository) Delete(id int) error {
	if _, err := r.db.Exec("DELETE FROM characters WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
