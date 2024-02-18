package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Excerpt struct {
	Id          int
	Title       string
	AuthorEmail string
	Created     time.Time
	Updated     time.Time
}

type ExcerptModel struct {
	Db *sql.DB
}

func (m ExcerptModel) Get(id int) (Excerpt, error) {
	stmt := `SELECT id, title, author_email, created, updated
	FROM excerpt WHERE id=?`

	var e Excerpt
	err := m.Db.QueryRow(stmt, id).Scan(&e.Id, &e.Title, &e.AuthorEmail,
		&e.Created, &e.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, errors.Join(ErrNoRecord, err)
		}
		return e, err
	}

	return e, nil
}

func (m ExcerptModel) Insert(title, author_email string) (int, error) {
	stmt := `INSERT INTO excerpt (title, author_email, created, updated)
	VALUES (?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	res, err := m.Db.Exec(stmt, title, author_email)
	if err != nil {
		var mysqlerr *mysql.MySQLError
		if errors.As(err, &mysqlerr) {
			if mysqlerr.Number == 1452 {
				return 0, ErrEmailDoesNotExist
			}
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m ExcerptModel) UpdateTitle(id int, title string) error {
	stmt := `UPDATE excerpt SET title=?, updated=UTC_TIMESTAMP() WHERE id=?`
	_, err := m.Db.Exec(stmt, title, id)
	if err != nil {
		return err
	}

	return nil
}

func (m ExcerptModel) GetByEmail(email string) ([]Excerpt, error) {
	stmt := `SELECT id, title, author_email, created, updated
	FROM excerpt
	WHERE author_email=?`
	rows, err := m.Db.Query(stmt, email)
	if err != nil {
		return []Excerpt{}, err
	}
	defer rows.Close()

	excerpts := []Excerpt{}
	for rows.Next() {
		e := Excerpt{}
		err = rows.Scan(&e.Id, &e.Title, &e.AuthorEmail, &e.Created, &e.Updated)
		if err != nil {
			return excerpts, err
		}
		excerpts = append(excerpts, e)
	}

	return excerpts, nil
}
