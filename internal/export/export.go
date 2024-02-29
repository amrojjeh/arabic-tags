package export

import (
	"slices"

	"github.com/amrojjeh/arabic-tags/internal/models"
)

type ExcerptExport struct {
	Version string
	Title   string
	Words   []WordExport
}

type WordExport struct {
	Word          string
	Connected     bool
	Punctuation   bool
	Ignore        bool
	SentenceStart bool
	Case          string
	State         string
}

func Export(e models.Excerpt, ws []models.Word) ExcerptExport {
	slices.SortFunc(ws, func(a, b models.Word) int {
		return a.WordPos - b.WordPos
	})

	export := ExcerptExport{
		Version: "0.1",
		Title:   e.Title,
		Words:   []WordExport{},
	}

	for _, w := range ws {
		export.Words = append(export.Words, WordExport{
			Word:          w.Word,
			Connected:     w.Connected,
			Punctuation:   w.Punctuation,
			Ignore:        w.Ignore,
			SentenceStart: w.SentenceStart,
			Case:          w.Case,
			State:         w.State,
		})
	}

	return export
}
