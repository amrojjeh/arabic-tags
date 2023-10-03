package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Grammar struct {
}

type Technical struct {
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

func (m ExcerptModel) Insert(title string) (uuid.UUID, error) {
	stmt := `INSERT INTO excerpt (id, title, content, grammar, technical,
	created, updated) VALUES (?, ?, "", "{}", "{}", UTC_TIMESTAMP(),
	UTC_TIMESTAMP())`

	// Technically it's not good practice to use UUIDv4 as a PK,
	// however, since we don't have an auth and we're using urls
	// to protect read and writes, the id needs to be hard to guess,
	// so we cannot just autoincrement it.

	// We're also unlikely to run into perf issues since we're operating on a
	// small scale.
	id := uuid.New()

	_, err := m.DB.Exec(stmt, [16]byte(id), title)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}
