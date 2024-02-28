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
	Id            int
	Word          string
	WordPos       int
	Connected     bool
	Punctuation   bool
	ExcerptId     int
	Ignore        bool
	SentenceStart bool
	Case          string
	State         string
	Created       time.Time
	Updated       time.Time
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
	stmt.WriteString(`INSERT INTO word (word, word_pos, connected, punctuation,
		excerpt_id, na_ignore, na_sentence_start, irab_case, irab_state, created, updated) VALUES `)
	vals := []any{}
	for i, w := range words {
		stmt.WriteString(`(?, ?, ?, ?, ?, false, false, "", "", UTC_TIMESTAMP(), UTC_TIMESTAMP())`)
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
	stmt := `SELECT id, word, word_pos, connected, punctuation, excerpt_id,
	na_ignore, na_sentence_start, irab_case, irab_state, created, updated
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
		err = rows.Scan(&w.Id, &w.Word, &w.WordPos, &w.Connected,
			&w.Punctuation, &w.ExcerptId, &w.Ignore, &w.SentenceStart,
			&w.Case, &w.State, &w.Created, &w.Updated)
		if err != nil {
			return nil, err
		}
		ws = append(ws, w)
	}

	return ws, nil
}

func (m WordModel) UpdateWord(id int, word string) error {
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

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos-1, updated=UTC_TIMESTAMP() WHERE id=?`, id)
	if err != nil {
		return err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos+1, updated=UTC_TIMESTAMP() WHERE id=?`, idOfOther)
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

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos+1, updated=UTC_TIMESTAMP() WHERE id=?`, id)
	if err != nil {
		return err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos-1, updated=UTC_TIMESTAMP() WHERE id=?`, idOfOther)
	if err != nil {
		return err
	}

	return t.Commit()
}

func (m WordModel) InsertAfter(id int, word string) (int, error) {
	t, err := m.Db.Begin()
	if err != nil {
		return 0, err
	}
	defer t.Rollback()
	var excerptId, wordPos int
	q := t.QueryRow(`SELECT excerpt_id, word_pos FROM word WHERE id=?`, id)
	err = q.Scan(&excerptId, &wordPos)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.Join(ErrNoRecord, err)
		}
		return 0, err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos+1, updated=UTC_TIMESTAMP() WHERE excerpt_id=? AND word_pos>?`,
		excerptId, wordPos)
	if err != nil {
		return 0, err
	}

	r, err := t.Exec(`INSERT INTO word (word, word_pos, connected, punctuation, excerpt_id, created, updated)
		VALUES (?, ?, false, false, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`,
		word, wordPos+1, excerptId)
	if err != nil {
		return 0, err
	}

	new_id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(new_id), t.Commit()
}

func (m WordModel) Delete(id int) error {
	t, err := m.Db.Begin()
	if err != nil {
		return err
	}
	defer t.Rollback()
	var excerptId, wordPos int
	q := t.QueryRow(`SELECT excerpt_id, word_pos FROM word WHERE id=?`, id)
	err = q.Scan(&excerptId, &wordPos)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.Join(ErrNoRecord, err)
		}
		return err
	}

	_, err = t.Exec(`UPDATE word SET word_pos=word_pos-1, updated=UTC_TIMESTAMP() WHERE excerpt_id=? AND word_pos>?`,
		excerptId, wordPos)
	if err != nil {
		return err
	}

	_, err = t.Exec(`DELETE FROM word WHERE id=?`, id)
	if err != nil {
		return err
	}

	return t.Commit()
}
func (m WordModel) Get(id int) (Word, error) {
	stmt := `SELECT id, word, word_pos, connected, excerpt_id, na_ignore,
	na_sentence_start, irab_case, irab_state, punctuation, created, updated
	FROM word
	WHERE id=?`

	var word Word
	q := m.Db.QueryRow(stmt, id)
	err := q.Scan(&word.Id, &word.Word, &word.WordPos, &word.Connected,
		&word.ExcerptId, &word.Ignore, &word.SentenceStart, &word.Case,
		&word.State, &word.Punctuation, &word.Created, &word.Updated)
	if err != nil {
		return Word{}, err
	}

	return word, nil
}

func (m WordModel) updateBool(id int, name string, val bool) error {
	stmt := fmt.Sprintf(`UPDATE word SET %v=?, UPDATED=UTC_TIMESTAMP()
	WHERE id=?`, name)

	_, err := m.Db.Exec(stmt, val, id)
	if err != nil {
		return err
	}

	return nil
}

func (m WordModel) UpdateConnect(id int, connect bool) error {
	return m.updateBool(id, "connected", connect)
}

func (m WordModel) UpdateIgnore(id int, ignore bool) error {
	return m.updateBool(id, "na_ignore", ignore)
}

func (m WordModel) UpdateSentenceStart(id int, sentence_start bool) error {
	return m.updateBool(id, "na_sentence_start", sentence_start)
}

func (m WordModel) UpdateIrab(id int, word_case, state string) error {
	stmt := `UPDATE word SET irab_case=?, irab_state=?, UPDATED=UTC_TIMESTAMP()
	WHERE id=?`

	_, err := m.Db.Exec(stmt, word_case, state, id)
	if err != nil {
		return err
	}

	return nil
}

func (m WordModel) UpdateState(id int, state string) error {
	stmt := `UPDATE word SET irab_state=?, UPDATED=UTC_TIMESTAMP()
	WHERE id=?`

	_, err := m.Db.Exec(stmt, state, id)
	if err != nil {
		return err
	}

	return nil
}
