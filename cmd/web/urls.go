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

func (u url) excerptEditSelectWord(id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v?word=%v", id, word_id)
}

func (u url) excerptEditWord(id any) string {
	return fmt.Sprintf("/excerpt/%v/word", id)
}

func (u url) excerptEditWordArgs(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word?word=%v", excerpt_id, word_id)
}

func (u url) excerptEditLetter(excerpt_id, word_id, letter_pos any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/letter/%v",
		excerpt_id, word_id, letter_pos)
}

func (u url) wordRight(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/right", excerpt_id, word_id)
}

func (u url) wordLeft(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/left", excerpt_id, word_id)
}

func (u url) wordAdd(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/add", excerpt_id, word_id)
}

func (u url) wordRemove(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/remove", excerpt_id, word_id)
}
