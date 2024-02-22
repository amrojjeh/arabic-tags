package main

import (
	"strconv"

	"github.com/amrojjeh/arabic-tags/internal/models"
	"github.com/amrojjeh/arabic-tags/ui/pages"
	"github.com/amrojjeh/arabic-tags/ui/partials"
	"github.com/amrojjeh/kalam"
	g "github.com/maragudk/gomponents"
)

func renderEdit(u url,
	e models.Excerpt, user models.User, ws []models.Word,
	selected int,
	error, warning string) g.Node {
	props := pages.EditProps{
		ExcerptTitle: e.Title,
		Username:     user.Username,
		Words:        []partials.WordProps{},
		Error:        error,
		Warning:      warning,
		TitleUrl:     u.excerptTitle(e.Id),
		EditWordUrl:  "",
		ExportUrl:    "#",
		LoginUrl:     u.login(),
		RegisterUrl:  u.register(),
		LogoutUrl:    u.logout(),
	}

	for _, w := range ws {
		wp := partials.WordProps{
			Id:          strconv.Itoa(w.Id),
			Word:        kalam.Prettify(w.Word),
			Punctuation: w.Punctuation,
			Connected:   w.Connected,
			Selected:    selected == w.WordPos,
			GetUrl:      u.excerptEditSelectWord(e.Id, w.WordPos),
		}
		props.Words = append(props.Words, wp)

		if w.WordPos == selected {
			props.SelectedWord.Word = w.Word
			props.SelectedWord.Id = strconv.Itoa(w.Id)
			props.SelectedWord.MoveRightUrl = u.wordRight(e.Id, w.Id)
			props.SelectedWord.MoveLeftUrl = u.wordLeft(e.Id, w.Id)
			props.EditWordUrl = u.excerptEditWordArgs(e.Id, w.WordPos)
			ls := kalam.LetterPacks(w.Word)
			for i, l := range ls {
				props.SelectedWord.Letters = append(props.SelectedWord.Letters,
					pages.LetterProps{
						Letter:          l.Unpointed(false),
						ShortVowel:      l.Vowel,
						Shadda:          l.Shadda,
						SuperscriptAlef: l.SuperscriptAlef,
						Index:           i,
						PostUrl:         u.excerptEditLetterArgs(e.Id, w.WordPos, i),
					})
			}
		}
	}

	return pages.EditPage(props)
}

func renderManuscript(u url,
	e models.Excerpt, ms models.Manuscript, user models.User,
	error, warning string) g.Node {
	props := pages.ManuscriptProps{
		ExcerptTitle:        e.Title,
		ReadOnly:            false,
		AcceptedPunctuation: kalam.PunctuationRegex().String(),
		Content:             ms.Content,
		Warning:             warning,
		Error:               error,
		Username:            user.Username,
		SaveUrl:             u.excerpt(e.Id),
		NextUrl:             u.excerptLock(e.Id),
		TitleUrl:            u.excerptTitle(e.Id),
		LoginUrl:            u.login(),
		RegisterUrl:         u.register(),
		LogoutUrl:           u.logout(),
	}

	return pages.ManuscriptPage(props)
}
