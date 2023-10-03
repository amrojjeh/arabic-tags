package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Grammar struct {
}

func (g *Grammar) Scan(src any) error {
	switch src.(type) {
	case []byte:
		json.Unmarshal(src.([]byte), g)
	default:
		return errors.New("grammar: cannot scan type")
	}

	return nil
}

type Technical struct {
}

func (t *Technical) Scan(src any) error {
	switch src.(type) {
	case []byte:
		json.Unmarshal(src.([]byte), t)
	default:
		return errors.New("technical: cannot scan type")
	}

	return nil
}

type Excerpt struct {
	ID        uuid.UUID
	Title     string
	Content   string
	Grammar   Grammar
	Technical Technical
	Created   time.Time
	Updated   time.Time
}

type ExcerptModel struct {
	DB *sql.DB
}

func (m ExcerptModel) Get(id uuid.UUID) (Excerpt, error) {
	stmt := `SELECT title, content, grammar, technical, created, updated
	FROM excerpt WHERE excerpt.id = UUID_TO_BIN(?)`

	var e Excerpt
	e.ID = id
	idVal, err := id.Value()
	if err != nil {
		return e, err
	}
	err = m.DB.QueryRow(stmt, idVal).Scan(&e.Title, &e.Content,
		&e.Grammar, &e.Technical, &e.Created, &e.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, ErrNoRecord
		}
		return e, err
	}

	return e, nil
}

func (m ExcerptModel) Insert(title string) (uuid.UUID, error) {
	stmt := `INSERT INTO excerpt (id, title, content, grammar, technical,
	created, updated) VALUES (UUID_TO_BIN(?), ?, "", "{}", "{}", UTC_TIMESTAMP(),
	UTC_TIMESTAMP())`

	// Technically it's not good practice to use UUIDv4 as a PK,
	// however, since we don't have an auth and we're using urls
	// to protect read and writes, the id needs to be hard to guess,
	// so we cannot just autoincrement it.

	// We're also unlikely to run into perf issues since we're operating on a
	// small scale.
	id := uuid.New()
	idVal, err := id.Value()
	if err != nil {
		return uuid.UUID{}, err
	}
	_, err = m.DB.Exec(stmt, idVal, title)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
