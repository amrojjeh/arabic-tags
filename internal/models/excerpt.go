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
	"github.com/google/uuid"
)

type Grammar struct {
	Words []GWord `json:"words"`
}

type GWord struct {
	Word     string   `json:"word"`
	Shrinked bool     `json:"shrinked"`
	LeftOver bool     `json:"leftOver"`
	Tags     []string `json:"tags"`

	// true if the word is preceding a punctuation or if the punctuation is
	// preceding a word (for rendering). Note that this is different from kalam's preceding
	Preceding   bool `json:"preceding"`
	Punctuation bool `json:"punctuation"`
}

func (g *Grammar) Scan(src any) error {
	switch src.(type) {
	case []byte:
		err := json.Unmarshal(src.([]byte), g)
		if err != nil {
			return err
		}
	default:
		return errors.New("grammar: cannot scan type")
	}

	return nil
}

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
	Id        int
	Title     string
	Content   string
	Grammar   Grammar
	Technical Technical
	CLocked   bool
	GLocked   bool
	Created   time.Time
	Updated   time.Time
}

type ExcerptModel struct {
	Db *sql.DB
}

// func (m ExcerptModel) Get(id uuid.UUID) (Excerpt, error) {
// 	stmt := `SELECT title, content, grammar, technical, c_locked, g_locked,
// 	c_share, g_share, t_share, created, updated
// 	FROM excerpt WHERE excerpt.id=UUID_TO_BIN(?)`

// 	var e Excerpt
// 	e.Id = id

// 	// UUID.Value always returns nil
// 	idVal, _ := id.Value()
// 	err := m.Db.QueryRow(stmt, idVal).Scan(&e.Title, &e.Content, &e.Grammar,
// 		&e.Technical, &e.CLocked, &e.GLocked, &e.Created, &e.Updated)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return e, ErrNoRecord
// 		}
// 		return e, err
// 	}

// 	return e, nil
// }

func (m ExcerptModel) Insert(title, author_email string) (int, error) {
	stmt := `INSERT INTO excerpt (title, author_id, created, updated)
	VALUES (?, (SELECT id FROM user WHERE email=?), UTC_TIMESTAMP(), UTC_TIMESTAMP())`

	res, err := m.Db.Exec(stmt, title, author_email)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m ExcerptModel) UpdateContent(id uuid.UUID, content string) error {
	stmt := `UPDATE excerpt SET content=?, updated=UTC_TIMESTAMP()
	WHERE id=UUID_TO_BIN(?) AND c_locked=FALSE`

	idVal, _ := id.Value()
	_, err := m.Db.Exec(stmt, content, idVal)
	if err != nil {
		return err
	}
	return nil
}

func (m ExcerptModel) SetContentLock(id uuid.UUID, lock bool) error {
	var stmt string
	if lock {
		stmt = `UPDATE excerpt SET c_locked=TRUE, updated=UTC_TIMESTAMP()
		WHERE id=UUID_TO_BIN(?)`
	} else {
		stmt = `UPDATE excerpt SET c_locked=FALSE, g_locked=FALSE, updated=UTC_TIMESTAMP()
		WHERE id=UUID_TO_BIN(?)`
	}

	idVal, _ := id.Value()
	_, err := m.Db.Exec(stmt, idVal)
	if err != nil {
		return err
	}
	return nil
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

func (m ExcerptModel) UpdateGrammar(id uuid.UUID, grammar Grammar) error {
	stmt := `UPDATE excerpt SET grammar=?, updated=UTC_TIMESTAMP()
	WHERE id=UUID_TO_BIN(?) AND g_locked=FALSE`

	idVal, _ := id.Value()
	load, err := json.Marshal(grammar)
	if err != nil {
		return err
	}
	_, err = m.Db.Exec(stmt, load, idVal)
	if err != nil {
		return err
	}

	return nil
}

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

func (e Excerpt) Export() ([]byte, error) {
	export := kalam.Excerpt{
		Name:      e.Title,
		Sentences: []kalam.Sentence{},
	}

	var sen *kalam.Sentence = nil
	for i, w := range e.Technical.Words {
		if w.SentenceStart {
			if sen != nil {
				export.Sentences = append(export.Sentences, *sen)
			}
			sen = &kalam.Sentence{
				Words: []kalam.Word{},
			}
		}
		sen.Words = append(sen.Words, kalam.Word{
			PointedWord: w.String(),
			Tags:        e.Grammar.Words[i].Tags,
			Punctuation: w.Punctuation,
			Ignore:      w.Ignore,
			Preceding:   w.Preceding,
		})
	}

	if sen != nil {
		export.Sentences = append(export.Sentences, *sen)
	}

	return json.Marshal(export)
}

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
	stmt := `SELECT e.id, e.title, e.created, e.updated
	FROM user AS u
	INNER JOIN excerpt AS e ON u.id = e.author_id
	WHERE u.email=?`
	rows, err := m.Db.Query(stmt, email)
	if err != nil {
		return []Excerpt{}, err
	}
	defer rows.Close()

	excerpts := []Excerpt{}
	for rows.Next() {
		e := Excerpt{}
		err = rows.Scan(&e.Id, &e.Title, &e.Created, &e.Updated)
		if err != nil {
			return excerpts, err
		}
		excerpts = append(excerpts, e)
	}

	return excerpts, nil
}
