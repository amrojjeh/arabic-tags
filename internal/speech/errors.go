package speech

import (
	"errors"
	"fmt"
)

var (
	ErrRequest     = errors.New("speech: request was not completed successfuly")
	ErrBadResponse = errors.New("speech: response buffer could not be read")
)

type BadFormatError struct {
	Text           string
	ExpectedFormat string
}

func (e BadFormatError) Error() string {
	return fmt.Sprintf("speech: text was not formatted properly (%v). Text=%v",
		e.ExpectedFormat, e.Text)
}

type UnrecognizedCharacterError struct {
	Character rune
}

func (e UnrecognizedCharacterError) Error() string {
	return fmt.Sprintf("speech: unrecognized character detected (0x%x)",
		e.Character)
}
