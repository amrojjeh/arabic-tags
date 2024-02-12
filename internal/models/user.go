package models

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	Db *sql.DB
}

type User struct {
	Id       int
	Username string
	Email    string
	Created  time.Time
	Updated  time.Time
}

func (m *UserModel) Register(username, email, password string) error {
	stmt := `INSERT INTO user (username, email, password_hash, created,
		updated) VALUES (?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	password_hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	_, err = m.Db.Exec(stmt, username, email, string(password_hash))
	if err != nil {
		return err
	}

	return nil
}

func (m *UserModel) Authenticate(email, password string) (bool, error) {
	stmt := `SELECT password_hash FROM user WHERE email=?`
	var pass_hash string
	res := m.Db.QueryRow(stmt, email)
	err := res.Scan(&pass_hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(pass_hash), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func (m *UserModel) Get(email string) (User, error) {
	stmt := `SELECT id, username, email, created, updated FROM user WHERE email=?`
	var user User
	res := m.Db.QueryRow(stmt, email)
	err := res.Scan(&user.Id, &user.Username, &user.Email, &user.Created, &user.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNoRecord
		}
		return User{}, err
	}
	return user, nil
}
