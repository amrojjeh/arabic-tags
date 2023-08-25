package speech

import "strings"

// TODO(Amr Ojjeh): Write documentation
// TODO(Amr Ojjeh): function to check if JSON is valid or corrupted

type Paragraph struct {
	Id        int        `json:"id"`
	Sentences []Sentence `json:"sentences"`
}

func (p Paragraph) String() string {
	sum := ""
	for _, s := range p.Sentences {
		sum += s.String() + "\n"
	}
	return sum
}

func (p *Paragraph) AddSentence(s *Sentence) {
	s.Id = len(p.Sentences)
	p.Sentences = append(p.Sentences, *s)
}

func (p *Paragraph) DeleteSentence(s Sentence) {
	cp := make([]Sentence, len(p.Sentences)-1, len(p.Sentences)+10)
	for i, v := range p.Sentences {
		if i > s.Id {
			cp[i-1] = v
			cp[i-1].Id = i - 1
		} else if i < s.Id {
			cp[i] = v
		}
	}
	p.Sentences = cp
}

func (p *Paragraph) DeleteSentenceId(id int64) {
	p.DeleteSentence(p.Sentences[id])
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

func (s *Sentence) AddWord(w Word) {
	w.Id = len(s.Words)
	s.Words = append(s.Words, w)
}

func NewSentence(v string) Sentence {
	ws := strings.Split(v, " ")
	words := make([]Word, len(ws))
	for i, w := range ws {
		words[i] = Word{
			Id:                 i,
			Value:              w,
			Case:               CaseNA,
			CaseIndicatorIndex: 0,
			CaseCause:          CauseNA}
	}
	return Sentence{Words: words}
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
