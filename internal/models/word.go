package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/amrojjeh/arabic-tags/internal/disambig"
)

type Word struct {
	Id          int
	Word        string
	WordPos     int
	Connected   bool
	Punctuation bool
	ExcerptId   int
	Created     time.Time
	Updated     time.Time
}

type WordModel struct {
	Db *sql.DB
}

func (m WordModel) DeleteByExcerptId(excerpt_id int) error {
	stmt := `DELETE FROM word WHERE excerpt_id=?`
	_, err := m.Db.Exec(stmt, excerpt_id)
	if err != nil {
		return err
	}

	return nil
}

func (m WordModel) GenerateWordsFromManuscript(ms Manuscript) error {
	words, err := disambig.Disambiguate(strings.TrimSpace(ms.Content))
	if err != nil {
		return err
	}
	if len(words) == 0 {
		return nil
	}

	var stmt strings.Builder
	stmt.WriteString(`INSERT INTO word (word, word_pos, connected, punctuation, excerpt_id, created, updated) VALUES `)
	vals := []any{}
	for i, w := range words {
		stmt.WriteString("(?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())")
		if i != len(words)-1 {
			stmt.WriteString(", ")
		}
		vals = append(vals, w.Word, i, w.Connected, w.Punctuation, ms.ExcerptId)
	}

	_, err = m.Db.Exec(stmt.String(), vals...)
	if err != nil {
		return errors.Join(fmt.Errorf("WordModel.GenerateWordsFromManuscript: could not execute stmt %v", stmt.String()), err)
	}

	return nil
}

func (m WordModel) GetWordsByExcerptId(excerpt_id int) ([]Word, error) {
	stmt := `SELECT id, word, word_pos, connected, punctuation, excerpt_id, created, updated
	FROM word
	WHERE excerpt_id=?
	ORDER BY word_pos`

	ws := []Word{}
	rows, err := m.Db.Query(stmt, excerpt_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var w Word
		err = rows.Scan(&w.Id, &w.Word, &w.WordPos, &w.Connected, &w.Punctuation, &w.ExcerptId, &w.Created, &w.Updated)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}

	return ws, nil
}

func (m WordModel) Update(id int, word string) error {
	stmt := `UPDATE word SET word=?, updated=UTC_TIMESTAMP() WHERE id=?`
	_, err := m.Db.Exec(stmt, word, id)
	if err != nil {
		return err
	}

	return nil
}

func (m WordModel) MoveRight(id int) error {
	t, err := m.Db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()

	var idOfOther int

	r := t.QueryRow(`SELECT w1.id
		FROM word AS w0
		INNER JOIN word AS w1
		ON w0.word_pos=w1.word_pos+1 AND w0.excerpt_id=w1.excerpt_id
		WHERE w0.id=? AND w0.word_pos>0`, id)

	err = r.Scan(&idOfOther)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ErrNotSwappable, err)
		}
		return err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos-1 WHERE id=?`, id)
	if err != nil {
		return err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos+1 WHERE id=?`, idOfOther)
	if err != nil {
		return err
	}

	return t.Commit()
}

func (m WordModel) MoveLeft(id int) error {
	t, err := m.Db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()

	var idOfOther int

	r := t.QueryRow(`SELECT w1.id
		FROM word AS w0
		INNER JOIN word AS w1
		ON w0.word_pos=w1.word_pos-1 AND w0.excerpt_id=w1.excerpt_id
		WHERE w0.id=? AND w1.word_pos>0`, id)

	err = r.Scan(&idOfOther)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ErrNotSwappable, err)
		}
		return err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos+1 WHERE id=?`, id)
	if err != nil {
		return err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos-1 WHERE id=?`, idOfOther)
	if err != nil {
		return err
	}

	return t.Commit()
}
