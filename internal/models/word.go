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
	Id        int
	Word      string
	WordPos   int
	ExcerptId int
	Created   time.Time
	Updated   time.Time
}

type WordModel struct {
	Db *sql.DB
}

func (m WordModel) GenerateWordsFromManuscript(ms Manuscript) error {
	words, err := disambig.Disambiguate(strings.TrimSpace(ms.Content))
	if err != nil {
		return err
	}
	if len(words) == 0 {
		return nil
	}

	vals := []any{}
	var stmt strings.Builder
	stmt.WriteString(`INSERT INTO word (word, word_pos, connected, excerpt_id, created, updated) VALUES `)
	for i, w := range words {
		stmt.WriteString("(?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())")
		if i != len(words)-1 {
			stmt.WriteString(", ")
		}
		vals = append(vals, w.Word, i, w.Connected, ms.ExcerptId)
	}

	_, err = m.Db.Exec(stmt.String(), vals...)
	if err != nil {
		return errors.Join(fmt.Errorf("WordModel.GenerateWordsFromManuscript: could not execute stmt %v", stmt.String()), err)
	}

	return nil
}
