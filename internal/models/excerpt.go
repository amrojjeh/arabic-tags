package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/amrojjeh/arabic-tags/internal/disambig"
	"github.com/amrojjeh/kalam"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Technical struct {
	Words []TWord `json:"words"`
}

func (t Technical) Text() string {
	par := ""
	for _, w := range t.Words {
		par += w.String()
		if !w.Preceding {
			par += " "
		}
	}
	par = strings.TrimSpace(par)
	return par
}

func (t Technical) TextWithoutPunctuation() string {
	par := ""
	for _, w := range t.Words {
		if w.Punctuation {
			par += " "
			continue
		}
		par += w.String()
		if !w.Preceding {
			par += " "
		}
	}
	return kalam.RemoveExtraWhitespace(par)
}

type TWord struct {
	Letters []Letter `json:"letters"`

	// Rendering data
	Preceding bool `json:"preceding"`

	// Word data
	Punctuation bool `json:"punctuation"`

	// Word data (configurable)
	SentenceStart bool `json:"sentenceStart"`
	Ignore        bool `json:"ignore"`
}

func (w TWord) String() string {
	if w.Punctuation {
		return w.Letters[0].Letter
	}
	word := ""
	for _, l := range w.Letters {
		word += l.String()
	}
	return word
}

type Letter struct {
	Letter string `json:"letter"`
	Vowel  string `json:"tashkeel"`
	Shadda bool   `json:"shadda"`
}

func (l Letter) String() string {
	if l.Shadda {
		return fmt.Sprintf("%v%v%c", l.Letter, l.Vowel, kalam.Shadda)
	}
	return fmt.Sprintf("%v%v", l.Letter, l.Vowel)
}

func (t *Technical) Scan(src any) error {
	switch src.(type) {
	case []byte:
		err := json.Unmarshal(src.([]byte), t)
		if err != nil {
			return err
		}
	default:
		return errors.New("technical: cannot scan type")
	}

	return nil
}

type Excerpt struct {
	Id          int
	Title       string
	AuthorEmail string
	Created     time.Time
	Updated     time.Time
}

type ExcerptModel struct {
	Db *sql.DB
}

func (m ExcerptModel) Get(id int) (Excerpt, error) {
	stmt := `SELECT id, title, author_email, created, updated
	FROM excerpt WHERE id=?`

	var e Excerpt
	err := m.Db.QueryRow(stmt, id).Scan(&e.Id, &e.Title, &e.AuthorEmail,
		&e.Created, &e.Updated)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e, errors.Join(ErrNoRecord, err)
		}
		return e, err
	}

	return e, nil
}

func (m ExcerptModel) Insert(title, author_email string) (int, error) {
	stmt := `INSERT INTO excerpt (title, author_email, created, updated)
	VALUES (?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	res, err := m.Db.Exec(stmt, title, author_email)
	if err != nil {
		var mysqlerr *mysql.MySQLError
		if errors.As(err, &mysqlerr) {
			if mysqlerr.Number == 1452 {
				return 0, ErrEmailDoesNotExist
			}
		}
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Assumes content is clean
// func (m ExcerptModel) ResetGrammar(id uuid.UUID) error {
// 	var generateWord = func(word string, punctuation bool, preceding bool) GWord {
// 		return GWord{
// 			Word:        word,
// 			Tags:        []string{},
// 			Punctuation: punctuation,
// 			Preceding:   preceding,
// 		}
// 	}
// 	excerpt, err := m.Get(id)
// 	if err != nil {
// 		return err
// 	}

// 	content := excerpt.Content
// 	words := make([]GWord, 0, len(strings.Split(content, " ")))
// 	word := ""
// 	for _, l := range content {
// 		if kalam.IsPunctuation(l) {
// 			if word != "" {
// 				words = append(words, generateWord(word, false, true))
// 				word = ""
// 			}
// 			// Assume preceding unless there's a space
// 			words = append(words, generateWord(string(l), true, true))
// 		} else if l == ' ' {
// 			wordCount := len(words)
// 			if word != "" {
// 				words = append(words, generateWord(word, false, false))
// 				word = ""
// 			} else if wordCount > 0 && words[wordCount-1].Punctuation {
// 				words[wordCount-1].Preceding = false
// 			}
// 		} else {
// 			word += string(l)
// 		}
// 	}
// 	if word != "" {
// 		words = append(words, generateWord(word, false, false))
// 	}
// 	grammar := Grammar{
// 		Words: words,
// 	}

// 	err = m.UpdateGrammar(id, grammar)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

func (m ExcerptModel) UpdateTechnical(id uuid.UUID, technical Technical) error {
	stmt := `UPDATE excerpt SET technical=?, updated=UTC_TIMESTAMP()
	WHERE excerpt.id=UUID_TO_BIN(?)`

	load, err := json.Marshal(technical)
	if err != nil {
		return err
	}

	idVal, _ := id.Value()
	_, err = m.Db.Exec(stmt, load, idVal)
	if err != nil {
		return err
	}

	return nil
}

// func (m ExcerptModel) ResetTechnical(id uuid.UUID) error {
// 	excerpt, err := m.Get(id)
// 	if err != nil {
// 		return err
// 	}

// 	technical := Technical{
// 		Words: make([]TWord, len(excerpt.Grammar.Words)),
// 	}
// 	for i, gw := range excerpt.Grammar.Words {
// 		technical.Words[i] = TWord{
// 			Letters:       make([]Letter, 0, utf8.RuneCountInString(gw.Word)),
// 			Punctuation:   gw.Punctuation,
// 			Preceding:     gw.Preceding || gw.Shrinked,
// 			SentenceStart: i == 0,
// 		}
// 		for _, l := range gw.Word {
// 			technical.Words[i].Letters = append(technical.Words[i].Letters, Letter{
// 				Letter: string(l),
// 				Vowel:  "",
// 				Shadda: false,
// 			})
// 		}
// 	}

// 	err = technical.Disambiguate()
// 	if err != nil {
// 		return err
// 	}
// 	err = m.UpdateTechnical(id, technical)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (e Excerpt) Export() ([]byte, error) {
// 	export := kalam.Excerpt{
// 		Name:      e.Title,
// 		Sentences: []kalam.Sentence{},
// 	}

// 	var sen *kalam.Sentence = nil
// 	for i, w := range e.Technical.Words {
// 		if w.SentenceStart {
// 			if sen != nil {
// 				export.Sentences = append(export.Sentences, *sen)
// 			}
// 			sen = &kalam.Sentence{
// 				Words: []kalam.Word{},
// 			}
// 		}
// 		sen.Words = append(sen.Words, kalam.Word{
// 			PointedWord: w.String(),
// 			Tags:        e.Grammar.Words[i].Tags,
// 			Punctuation: w.Punctuation,
// 			Ignore:      w.Ignore,
// 			Preceding:   w.Preceding,
// 		})
// 	}

// 	if sen != nil {
// 		export.Sentences = append(export.Sentences, *sen)
// 	}

// 	return json.Marshal(export)
// }

// TODO(Amr Ojjeh): Automatically vowelize mabni words
func (t *Technical) Disambiguate() error {
	dWords, err := disambig.Disambiguate(t.TextWithoutPunctuation())
	if err != nil {
		return err
	}

	mapper := struct {
		li int // letter index
		wi int // word index
	}{li: -1, wi: 0}
	for _, dWord := range dWords {
		for _, dLetter := range dWord {
			for t.Words[mapper.wi].Punctuation {
				mapper.li = 0
				mapper.wi += 1
			}
			mapper.li += 1
			if mapper.li == len(t.Words[mapper.wi].Letters) {
				mapper.li = 0
				mapper.wi += 1
				for t.Words[mapper.wi].Punctuation {
					mapper.wi += 1
				}
			}
			letter := &t.Words[mapper.wi].Letters[mapper.li]
			if dLetter.Vowel != 0 {
				letter.Vowel = string(dLetter.Vowel)
			} else {
				letter.Vowel = string(kalam.Sukoon)
			}
			letter.Shadda = dLetter.Shadda
		}
	}
	return nil
}

func (m *ExcerptModel) GetByEmail(email string) ([]Excerpt, error) {
	stmt := `SELECT id, title, author_email, created, updated
	FROM excerpt
	WHERE author_email=?`
	rows, err := m.Db.Query(stmt, email)
	if err != nil {
		return []Excerpt{}, err
	}
	defer rows.Close()

	excerpts := []Excerpt{}
	for rows.Next() {
		e := Excerpt{}
		err = rows.Scan(&e.Id, &e.Title, &e.AuthorEmail, &e.Created, &e.Updated)
		if err != nil {
			return excerpts, err
		}
		excerpts = append(excerpts, e)
	}

	return excerpts, nil
}
