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
