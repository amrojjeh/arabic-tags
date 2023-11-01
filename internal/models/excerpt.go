package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Grammar struct {
	Words []Word `json:"words"`
}

type Word struct {
	Word     string   `json:"word"`
	Shrinked bool     `json:"shrinked"`
	LeftOver bool     `json:"leftOver"`
	Tags     []string `json:"tags"`
}

func (g *Grammar) Scan(src any) error {
	switch src.(type) {
	case []byte:
		err := json.Unmarshal(src.([]byte), g)
		if err != nil {
			return err
		}
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
		err := json.Unmarshal(src.([]byte), t)
		if err != nil {
			return err
		}
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
	CLocked   bool
	GLocked   bool
	CShare    uuid.UUID
	GShare    uuid.UUID
	TShare    uuid.UUID
	Created   time.Time
	Updated   time.Time
}

type ExcerptModel struct {
	DB *sql.DB
}

func (m ExcerptModel) Get(id uuid.UUID) (Excerpt, error) {
	stmt := `SELECT title, content, grammar, technical, c_locked, g_locked,
	c_share, g_share, t_share, created, updated
	FROM excerpt WHERE excerpt.id=UUID_TO_BIN(?)`

	var e Excerpt
	e.ID = id

	// UUID.Value always returns nil
	idVal, _ := id.Value()
	err := m.DB.QueryRow(stmt, idVal).Scan(&e.Title, &e.Content, &e.Grammar,
		&e.Technical, &e.CLocked, &e.GLocked, &e.CShare, &e.GShare, &e.TShare,
		&e.Created, &e.Updated)
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
	c_locked, g_locked, c_share, g_share, t_share, created, updated)
	VALUES (UUID_TO_BIN(?), ?, "", "{}", "{}", FALSE, FALSE, UUID_TO_BIN(?),
	UUID_TO_BIN(?), UUID_TO_BIN(?), UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	// Technically it's not good practice to use UUIDv4 as a PK,
	// however, since we don't have an auth and we're using urls
	// to protect read and writes, the id needs to be hard to guess,
	// so we cannot just autoincrement it.

	// We're also unlikely to run into perf issues since we're operating on a
	// small scale.
	id := uuid.New()
	idVal, _ := id.Value()

	cShareVal, _ := uuid.New().Value()
	gShareVal, _ := uuid.New().Value()
	tShareVal, _ := uuid.New().Value()
	_, err := m.DB.Exec(stmt, idVal, title, cShareVal, gShareVal, tShareVal)
	if err != nil {
		return uuid.UUID{}, err
	}
	return id, nil
}

func (m ExcerptModel) UpdateContent(id uuid.UUID, content string) error {
	stmt := `UPDATE excerpt SET content=?, updated=UTC_TIMESTAMP()
	WHERE id=UUID_TO_BIN(?) AND c_locked=FALSE`

	idVal, _ := id.Value()
	_, err := m.DB.Exec(stmt, content, idVal)
	if err != nil {
		return err
	}
	return nil
}

func (m ExcerptModel) SetContentLock(id uuid.UUID, lock bool) error {
	stmt := `UPDATE excerpt SET c_locked=?, updated=UTC_TIMESTAMP()
	WHERE id=UUID_TO_BIN(?)`

	idVal, _ := id.Value()
	_, err := m.DB.Exec(stmt, lock, idVal)
	if err != nil {
		return err
	}
	return nil
}

func (m ExcerptModel) ResetGrammar(id uuid.UUID) error {
	// TODO(Amr Ojjeh): Add g_locked
	excerpt, err := m.Get(id)
	if err != nil {
		return err
	}

	content := excerpt.Content

	strWords := strings.Split(content, " ")
	words := make([]Word, len(strWords))
	for i, s := range strWords {
		words[i] = Word{
			Word:     s,
			Shrinked: false,
			LeftOver: false,
			Tags:     []string{},
		}
	}
	grammar := Grammar{
		Words: words,
	}

	err = m.UpdateGrammar(id, grammar)
	if err != nil {
		return err
	}

	return nil
}

func (m ExcerptModel) UpdateGrammar(id uuid.UUID, grammar Grammar) error {
	stmt := `UPDATE excerpt SET grammar=?, updated=UTC_TIMESTAMP()
	WHERE id=UUID_TO_BIN(?)`

	idVal, _ := id.Value()
	load, err := json.Marshal(grammar)
	if err != nil {
		return err
	}
	_, err = m.DB.Exec(stmt, load, idVal)
	if err != nil {
		return err
	}

	return nil
}
