package models

import "errors"

var (
	ErrNoRecord          = errors.New("models: no matching record found")
	ErrDuplicateUsername = errors.New("models: username already taken")
	ErrDuplicateEmail    = errors.New("models: email already taken")
	ErrEmailDoesNotExist = errors.New("models: email does not exist")
	ErrNotSwappable      = errors.New("models: word can't be swapped")
	ErrInvalidValue      = errors.New("models: invalid value")
)
