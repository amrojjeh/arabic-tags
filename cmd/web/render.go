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
	selectedId int,
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
			Selected:    selectedId == w.Id,
			GetUrl:      u.excerptEditSelectWord(e.Id, w.Id),
		}
		props.Words = append(props.Words, wp)

		if w.Id == selectedId {
			props.SelectedWord = pages.SelectedWordProps{
				Id:               strconv.Itoa(w.Id),
				Word:             w.Word,
				Letters:          []pages.LetterProps{},
				Connected:        w.Connected,
				Ignore:           w.Ignore,
				SentenceStart:    w.SentenceStart,
				MoveRightUrl:     u.wordRight(e.Id, w.Id),
				MoveLeftUrl:      u.wordLeft(e.Id, w.Id),
				AddWordUrl:       u.wordAdd(e.Id, w.Id),
				RemoveWordUrl:    u.wordRemove(e.Id, w.Id),
				ConnectedUrl:     u.wordConnect(e.Id, w.Id),
				IgnoreUrl:        u.wordIgnore(e.Id, w.Id),
				SentenceStartUrl: u.wordSentenceStart(e.Id, w.Id),
			}
			props.EditWordUrl = u.excerptEditWordArgs(e.Id, w.Id)
			ls := kalam.LetterPacks(w.Word)
			for i, l := range ls {
				props.SelectedWord.Letters = append(props.SelectedWord.Letters,
					pages.LetterProps{
						Letter:          l.Unpointed(false),
						ShortVowel:      l.Vowel,
						Shadda:          l.Shadda,
						SuperscriptAlef: l.SuperscriptAlef,
						Index:           i,
						PostUrl:         u.excerptEditLetter(e.Id, w.Id, i),
					})
			}
		}
	}

	return pages.EditPage(props)
}

func renderText(u url,
	e models.Excerpt, words []models.Word,
	selectedId int) g.Node {
	wps := []partials.WordProps{}
	for _, word := range words {
		wps = append(wps, partials.WordProps{
			Id:          strconv.Itoa(word.Id),
			Word:        kalam.Prettify(word.Word),
			Punctuation: word.Punctuation,
			Connected:   word.Connected,
			Selected:    word.Id == selectedId,
			GetUrl:      u.excerptEditSelectWord(e.Id, word.Id),
		})
	}

	return partials.Text(wps)
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
