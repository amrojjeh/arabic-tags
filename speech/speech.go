package speech

import "strings"

// TODO(Amr Ojjeh): Write documentation

type Sentence struct {
	Words []Word
}

func (s Sentence) String() string {
	sum := ""
	for _, w := range s.Words {
		sum += w.Value + " "
	}
	return sum
}

func NewSentence(v string) Sentence {
	ws := strings.Split(" ", v)
	words := make([]Word, len(ws))
	for i, w := range ws {
		words[i] = Word{
			Value:              w,
			Case:               CaseNA,
			CaseIndicatorIndex: 0,
			CaseCause:          CauseNA}
	}
	return Sentence{Words: words}
}

type Word struct {
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
