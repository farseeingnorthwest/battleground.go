package storage

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/farseeingnorthwest/playground/battlefield/v2"
	"github.com/jmoiron/sqlx"
)

type Skill struct {
	ID      int
	Name    string
	Reactor *Reactor
}

type Reactor battlefield.FatReactor

func (r *Reactor) Value() (driver.Value, error) {
	return json.Marshal((*battlefield.FatReactor)(r))
}

func (r *Reactor) Scan(value interface{}) error {
	j, ok := value.([]byte)
	if !ok {
		return errors.New("invalid argument")
	}

	var f battlefield.FatReactorFile
	if err := json.Unmarshal(j, &f); err != nil {
		return err
	}

	*r = Reactor(*f.FatReactor)
	return nil
}

type SkillRepository struct {
	db *sqlx.DB
}

func NewSkillRepository(db *sqlx.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

type SkillQuery struct {
	order  string
	limit  string
	offset string
}

func (q *SkillQuery) OrderBy(asc bool) *SkillQuery {
	if asc {
		q.order = "ORDER BY id ASC"
	} else {
		q.order = "ORDER BY id DESC"
	}

	return q
}

func (q *SkillQuery) Limit(limit int) *SkillQuery {
	q.limit = "LIMIT " + strconv.Itoa(limit)
	return q
}

func (q *SkillQuery) Offset(offset int) *SkillQuery {
	q.offset = "OFFSET " + strconv.Itoa(offset)
	return q
}

func (q *SkillQuery) Select() (string, []any) {
	qs := strings.Join(
		[]string{
			"SELECT * FROM skills ORDER BY id",
			q.order,
			q.limit,
			q.offset,
		},
		" ",
	)

	return qs, nil
}

func (r SkillRepository) Find(query *SkillQuery) (skills []Skill, err error) {
	qs, args := query.Select()
	err = r.db.Select(&skills, qs, args...)
	return
}

func (r SkillRepository) Get(id int) (*Skill, error) {
	var skill Skill
	if err := r.db.Get(&skill, "SELECT * FROM skills WHERE id = $1", id); err != nil {
		return nil, err
	}

	return &skill, nil
}

func (r SkillRepository) Create(skill *Skill) error {
	if err := r.db.Get(skill, "INSERT INTO skills (name, reactor) VALUES ($1, $2) RETURNING *", skill.Name, skill.Reactor); err != nil {
		return err
	}

	return nil
}

func (r SkillRepository) Update(skill *Skill) error {
	if err := r.db.Get(skill, "UPDATE skills SET name = $1, reactor = $2 WHERE id = $3 RETURNING *", skill.Name, skill.Reactor, skill.ID); err != nil {
		return err
	}

	return nil
}

func (r SkillRepository) Delete(id int) error {
	if _, err := r.db.Exec("DELETE FROM skills WHERE id = $1", id); err != nil {
		return err
	}

	return nil
}
