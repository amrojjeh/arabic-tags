package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/amrojjeh/arabic-tags/internal/speech"
	"github.com/google/uuid"
)

type Grammar struct {
	Words []GWord `json:"words"`
}

type GWord struct {
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
	Words []TWord
}

type TWord struct {
	Letters []Letter `json:"letters"`
	Tags    []string `json:"tags"`
}

func (w TWord) String() string {
	word := ""
	for _, l := range w.Letters {
		word += l.String()
	}
	return word
}

type Letter struct {
	Letter   string `json:"letter"`
	Tashkeel string `json:"tashkeel"`
	Shadda   bool   `json:"shadda"`
}

func (l Letter) String() string {
	if l.Shadda {
		return fmt.Sprintf("%v%v%v", l.Letter, l.Tashkeel, speech.Shadda)
	}
	return fmt.Sprintf("%v%v", l.Letter, l.Tashkeel)
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

func (m ExcerptModel) GetSharedContent(cShare uuid.UUID) (Excerpt, error) {
	stmt := `SELECT id, title, content, grammar, technical, c_locked, g_locked,
	g_share, t_share, created, updated
	FROM excerpt WHERE excerpt.c_share=UUID_TO_BIN(?)`

	var e Excerpt
	e.CShare = cShare

	cShareVal, _ := cShare.Value()
	err := m.DB.QueryRow(stmt, cShareVal).Scan(&e.ID, &e.Title, &e.Content,
		&e.Grammar, &e.Technical, &e.CLocked, &e.GLocked, &e.GShare, &e.TShare,
		&e.Created, &e.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, ErrNoRecord
		}
		return e, err
	}
	return e, nil
}

func (m ExcerptModel) GetSharedGrammar(gShare uuid.UUID) (Excerpt, error) {
	stmt := `SELECT id, title, content, grammar, technical, c_locked, g_locked,
	c_share, t_share, created, updated
	FROM excerpt WHERE excerpt.g_share=UUID_TO_BIN(?)`

	var e Excerpt
	e.GShare = gShare

	gShareVal, _ := gShare.Value()
	err := m.DB.QueryRow(stmt, gShareVal).Scan(&e.ID, &e.Title, &e.Content,
		&e.Grammar, &e.Technical, &e.CLocked, &e.GLocked, &e.CShare, &e.TShare,
		&e.Created, &e.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, ErrNoRecord
		}
		return e, err
	}
	return e, nil
}

func (m ExcerptModel) GetSharedTechnical(tShare uuid.UUID) (Excerpt, error) {
	stmt := `SELECT id, title, content, grammar, technical, c_locked, g_locked,
	c_share, g_share, created, updated
	FROM excerpt WHERE excerpt.t_share=UUID_TO_BIN(?)`

	var e Excerpt
	e.TShare = tShare

	tShareVal, _ := tShare.Value()
	err := m.DB.QueryRow(stmt, tShareVal).Scan(&e.ID, &e.Title, &e.Content,
		&e.Grammar, &e.Technical, &e.CLocked, &e.GLocked, &e.CShare, &e.GShare,
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

func (m ExcerptModel) UpdateSharedContent(cShare uuid.UUID, content string) error {
	stmt := `UPDATE excerpt SET content=?, updated=UTC_TIMESTAMP()
	WHERE c_share=UUID_TO_BIN(?) AND c_locked=FALSE`

	idVal, _ := cShare.Value()
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

func (m ExcerptModel) SetGrammarLock(id uuid.UUID, lock bool) error {
	stmt := `UPDATE excerpt SET g_locked=?, updated=UTC_TIMESTAMP()
	WHERE id=UUID_TO_BIN(?)`

	idVal, _ := id.Value()
	_, err := m.DB.Exec(stmt, lock, idVal)
	if err != nil {
		return err
	}

	return nil
}

func (m ExcerptModel) ResetGrammar(id uuid.UUID) error {
	excerpt, err := m.Get(id)
	if err != nil {
		return err
	}

	content := excerpt.Content

	strWords := strings.Split(content, " ")
	words := make([]GWord, len(strWords))
	for i, s := range strWords {
		words[i] = GWord{
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
	WHERE id=UUID_TO_BIN(?) AND g_locked=FALSE`

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

func (m ExcerptModel) UpdateSharedGrammar(gShare uuid.UUID, grammar Grammar) error {
	stmt := `UPDATE excerpt SET grammar=?, updated=UTC_TIMESTAMP()
	WHERE g_share=UUID_TO_BIN(?) AND g_locked=FALSE`

	idVal, _ := gShare.Value()
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

func (m ExcerptModel) UpdateTechnical(id uuid.UUID, technical Technical) error {
	stmt := `UPDATE excerpt SET technical=?, updated=UTC_TIMESTAMP()
	WHERE excerpt.id=UUID_TO_BIN(?)`

	load, err := json.Marshal(technical)
	if err != nil {
		return err
	}

	idVal, _ := id.Value()
	_, err = m.DB.Exec(stmt, load, idVal)
	if err != nil {
		return err
	}

	return nil
}

func (m ExcerptModel) ResetTechnical(id uuid.UUID) error {
	excerpt, err := m.Get(id)
	if err != nil {
		return err
	}

	technical := Technical{
		Words: make([]TWord, len(excerpt.Grammar.Words)),
	}
	for i, gw := range excerpt.Grammar.Words {
		technical.Words[i] = TWord{
			Letters: make([]Letter, len(gw.Word)),
			Tags:    []string{},
		}
		for x, l := range gw.Word {
			technical.Words[i].Letters[x] = Letter{
				Letter:   string(l),
				Tashkeel: "",
				Shadda:   false,
			}
		}
	}

	err = m.UpdateTechnical(id, technical)
	if err != nil {
		return err
	}
	return nil
}
