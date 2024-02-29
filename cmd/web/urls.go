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

func (u url) wordSelect(id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v?word=%v", id, word_id)
}

func (u url) wordEdit(id any) string {
	return fmt.Sprintf("/excerpt/%v/word", id)
}

func (u url) wordEditArgs(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word?word=%v", excerpt_id, word_id)
}

func (u url) letterEdit(excerpt_id, word_id, letter_pos any) string {
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

func (u url) wordConnect(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/connect", excerpt_id, word_id)
}

func (u url) wordSentenceStart(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/sentence_start", excerpt_id, word_id)
}

func (u url) wordIgnore(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/ignore", excerpt_id, word_id)
}

func (u url) wordCase(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/case", excerpt_id, word_id)
}

func (u url) wordState(excerpt_id, word_id any) string {
	return fmt.Sprintf("/excerpt/%v/word/%v/state", excerpt_id, word_id)
}
