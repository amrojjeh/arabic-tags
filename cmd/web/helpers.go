package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

func (app *application) clientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	app.logger.Error("server error", slog.String("error", err.Error()))
	http.Error(w, http.StatusText(http.StatusInternalServerError),
		http.StatusInternalServerError)
}

func (app *application) noBody(w http.ResponseWriter) {
	w.Write([]byte(""))
}

func cleanContent(content string) (string, error) {
	for _, c := range content {
		if !(isArabicLetter(c) || isWhitespace(c)) {
			return "", errors.New(fmt.Sprintf("%v is an invalid letter", c))
		}
	}

	// Remove double spaces
	r, _ := regexp.Compile(" +")
	content = r.ReplaceAllString(content, " ")

	// Trim sentence
	content = strings.TrimFunc(content, unicode.IsSpace)
	return content, nil
}

// isArabicLetter does not include tashkeel
func isArabicLetter(letter rune) bool {
	if letter >= 0x0621 && letter <= 0x063A {
		return true
	}
	if letter >= 0x0641 && letter <= 0x064A {
		return true
	}
	return false
}

func isWhitespace(letter rune) bool {
	return letter == ' '
}
