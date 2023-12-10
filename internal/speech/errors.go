package speech

import (
	"errors"
	"fmt"
)

var (
	ErrRequest               = errors.New("speech: request was not completed successfuly")
	ErrBadResponse           = errors.New("speech: response buffer could not be read")
	ErrUnrecognizedCharacter = errors.New("")
)

type UnrecognizedCharacterError struct {
	Character rune
}

func (e UnrecognizedCharacterError) Error() string {
	return fmt.Sprintf("speech: unrecognized character detected (0x%x)",
		e.Character)
}
