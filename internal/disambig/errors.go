package disambig

import (
	"errors"
	"fmt"
)

var (
	ErrRequest     = errors.New("disambiguate: request was not completed successfuly")
	ErrBadResponse = errors.New("disambiguate: response buffer could not be read")
	ErrTruncated   = errors.New("disambiguate: response truncated input")
)

type BadFormatError struct {
	Text           string
	ExpectedFormat string
}

func (e BadFormatError) Error() string {
	return fmt.Sprintf("disambiguate: text was not formatted properly (%v). Text=%v",
		e.ExpectedFormat, e.Text)
}

type UnrecognizedCharacterError struct {
	Character rune
}

func (e UnrecognizedCharacterError) Error() string {
	return fmt.Sprintf("disambiguate: unrecognized character detected (0x%x)",
		e.Character)
}
