package speech

import (
	"errors"
	"fmt"
	"strings"
)

// TODO(Amr Ojjeh): Write documentation
// TODO(Amr Ojjeh): function to check if JSON is valid or corrupted

func wordsSlice(id int, value string) (int, []Word) {
	ws := strings.Split(value, " ")
	words := make([]Word, len(ws))
	for i, w := range ws {
		words[i] = Word{
			Id:                 id,
			Value:              w,
			Case:               CaseNA,
			CaseIndicatorIndex: 0,
			CaseCause:          CauseNA}
		id++
	}
	return id, words
}

type Paragraph struct {
	AvailableId int        `json:"available_id"`
	Sentences   []Sentence `json:"sentences"`
}

func (p Paragraph) String() string {
	sum := ""
	for _, s := range p.Sentences {
		sum += s.String() + "\n"
	}
	return sum
}

func (p *Paragraph) EditSentence(id int, value string) bool {
	sen, err := p.GetSentenceId(id)
	if err != nil {
		return false
	}
	newId, words := wordsSlice(p.AvailableId, value)
	p.AvailableId = newId
	sen.Words = words
	return true
}

func (p *Paragraph) AddSentence(value string) Sentence {
	id, words := wordsSlice(p.AvailableId, value)
	p.AvailableId = id
	s := Sentence{Id: p.AvailableId,
		Words: words}
	p.AvailableId++
	p.Sentences = append(p.Sentences, s)
	return s
}

// DeleteSentence deletes a sentence. If a sentence was deleted, then it returns true. Otherwise, it returns false
func (p *Paragraph) DeleteSentence(id int) bool {
	si, err := p.GetSentenceIndex(id)
	if err != nil {
		return false
	}
	cp := make([]Sentence, len(p.Sentences)-1, len(p.Sentences)+10)
	for i, v := range p.Sentences {
		if i > si {
			cp[i-1] = v
		} else if i < si {
			cp[i] = v
		}
	}
	p.Sentences = cp
	return true
}

func (p *Paragraph) GetSentenceIndex(id int) (int, error) {
	for i, v := range p.Sentences {
		if v.Id == id {
			return i, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("sentence id with %v was not found", id))
}

func (p *Paragraph) swapSentence(o1, o2 int) {
	tempS := p.Sentences[o1]
	p.Sentences[o1] = p.Sentences[o2]
	p.Sentences[o2] = tempS
}

func (p *Paragraph) MoveSentenceUp(id int) {
	if s1, err := p.GetSentenceIndex(id); err == nil {
		if s2 := s1 - 1; s1 > 0 {
			p.swapSentence(s1, s2)
		}
	}
}

func (p *Paragraph) MoveSentenceDown(id int) {
	if s1, err := p.GetSentenceIndex(id); err == nil {
		if s2 := s1 + 1; s2 < len(p.Sentences) {
			p.swapSentence(s1, s2)
		}
	}
}

func (p *Paragraph) GetSentenceId(id int) (*Sentence, error) {
	si, err := p.GetSentenceIndex(id)
	if err != nil {
		return &Sentence{}, errors.New(fmt.Sprintf("sentence id with %v was not found", id))
	}
	return &p.Sentences[si], nil
}

type Sentence struct {
	Id    int    `'json:"id"`
	Words []Word `json:"words"`
}

func (s Sentence) String() string {
	sum := ""
	for _, w := range s.Words {
		sum += w.Value + " "
	}
	return sum
}

// TODO(Amr Ojjeh): Save JSON in SafeBW
type Word struct {
	Id                 int       `json:"id"`
	Value              string    `json:"value"`
	Case               caseClass `json:"case"`
	CaseIndicatorIndex int       `json:"case_index"`
	CaseCause          caseCause `json:"case_cause"`
}

type caseClass string

const (
	CaseNA         caseClass = "CASE_NA"
	CaseNominative caseClass = "CASE_NOMINATIVE"
	CaseAccusative caseClass = "CASE_ACCUSATIVE"
	CaseGenitive   caseClass = "CASE_GENITIVE"
	CaseJussive    caseClass = "CASE_JUSSIVE"
)

type caseCause string

const (
	CauseNA        caseCause = "CAUSE_NA"
	CausePredicate caseCause = "CAUSE_PREDICATE"
	CauseSubject   caseCause = "CAUSE_SUBJECT"
	CausePastVerb  caseCause = "CAUSE_PAST_VERB"
)
