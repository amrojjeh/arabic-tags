package models

import (
	"database/sql"
	"errors"
	"time"
)

type ManuscriptModel struct {
	Db *sql.DB
}

type Manuscript struct {
	Id        int
	Content   string
	ExcerptId int
	Created   time.Time
	Updated   time.Time
}

func (m ManuscriptModel) Insert(excerpt_id int) (int, error) {
	stmt := `INSERT INTO manuscript (content, excerpt_id, created,
	updated) VALUES ("", ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	res, err := m.Db.Exec(stmt, excerpt_id)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m ManuscriptModel) Update(id int, content string) error {
	stmt := `UPDATE manuscript
	SET content=?, updated=UTC_TIMESTAMP()
	WHERE id=?`

	_, err := m.Db.Exec(stmt, content, id)
	if err != nil {
		return err
	}
	return nil
}

func (m ManuscriptModel) UpdateByExcerptId(excerpt_id int, content string) error {
	stmt := `UPDATE manuscript
	SET content=?, updated=UTC_TIMESTAMP()
	WHERE excerpt_id=?`

	_, err := m.Db.Exec(stmt, content, excerpt_id)
	if err != nil {
		return err
	}
	return nil
}

func (m ManuscriptModel) Get(id int) (Manuscript, error) {
	stmt := `SELECT id, content, excerpt_id, created, updated
	FROM manuscript
	WHERE id=?`

	row := m.Db.QueryRow(stmt, id)
	ms := Manuscript{}
	err := row.Scan(&ms.Id, &ms.Content, &ms.ExcerptId,
		&ms.Created, &ms.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Manuscript{}, errors.Join(ErrNoRecord, err)
		}
		return Manuscript{}, err
	}
	return ms, nil
}

func (m ManuscriptModel) GetByExcerptId(excerpt_id int) (Manuscript, error) {
	stmt := `SELECT id, content, excerpt_id, created, updated
	FROM manuscript
	WHERE excerpt_id=?`

	row := m.Db.QueryRow(stmt, excerpt_id)
	ms := Manuscript{}
	err := row.Scan(&ms.Id, &ms.Content, &ms.ExcerptId,
		&ms.Created, &ms.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Manuscript{}, errors.Join(ErrNoRecord, err)
		}
		return Manuscript{}, err
	}
	return ms, nil
}

func (m ManuscriptModel) Delete(excerpt_id int) error {
	stmt := `DELETE FROM manuscript WHERE excerpt_id=?`
	_, err := m.Db.Exec(stmt, excerpt_id)
	if err != nil {
		return err
	}
	return nil
}
