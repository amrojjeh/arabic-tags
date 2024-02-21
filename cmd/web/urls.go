package main

import "fmt"

type url struct{}

func (u url) index() string {
	return "/"
}

func (u url) login() string {
	return "/login"
}

func (u url) register() string {
	return "/register"
}

func (u url) logout() string {
	return "/logout"
}

func (u url) home() string {
	return "/home"
}

func (u url) createExcerpt() string {
	return "/excerpt"
}

func (u url) excerpt(id any) string {
	return fmt.Sprintf("/excerpt/%v", id)
}

func (u url) excerptTitle(id any) string {
	return fmt.Sprintf("/excerpt/%v/title", id)
}

func (u url) excerptLock(id any) string {
	return fmt.Sprintf("/excerpt/%v/lock", id)
}

func (u url) excerptEditSelectWord(id, word_pos any) string {
	return fmt.Sprintf("/excerpt/%v?word_pos=%v", id, word_pos)
}

func (u url) excerptEditWord(id any) string {
	return fmt.Sprintf("/excerpt/%v/word", id)
}

func (u url) excerptEditWordArgs(id, word_pos any) string {
	return fmt.Sprintf("/excerpt/%v/word?word_pos=%v", id, word_pos)
}

func (u url) excerptEditLetter(id any) string {
	return fmt.Sprintf("/excerpt/%v/letter", id)
}

func (u url) excerptEditLetterArgs(id, word_pos, letter_pos any) string {
	return fmt.Sprintf("/excerpt/%v/letter?word_pos=%v&letter_pos=%v",
		id, word_pos, letter_pos)
}
