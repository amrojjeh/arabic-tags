package models

import "errors"

var (
	ErrNoRecord          = errors.New("models: no matching record found")
	ErrDuplicateUsername = errors.New("models: username already taken")
	ErrDuplicateEmail    = errors.New("models: email already taken")
)
