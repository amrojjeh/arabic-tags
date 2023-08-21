package speech

// TODO(Amr Ojjeh): Write documentation

type Sentence []Word

func (s Sentence) String() string {
	sum := ""
	for _, w := range s {
		sum += w.String() + " "
	}
	return sum
}

type Word []Token

func (w Word) String() string {
	sum := ""
	for _, t := range w {
		sum += t.Value
	}
	return sum
}

type Token struct {
	Value string
	Case  CaseType
}

type CaseType struct {
	Type      caseClass
	Indicator caseIndicator
}

type caseIndicator string

const (
	IndicatorNA       caseIndicator = "INDICATOR_NA"
	IndicatorImplicit caseIndicator = "INDICATOR_IMPLICIT"
	IndicatorDammah   caseIndicator = "INDICATOR_DAMMAH"
)

type caseClass string

const (
	CaseNA         caseClass = "CASE_NA"
	CaseNominative caseClass = "CASE_NOMINATIVE"
	CaseAccusative caseClass = "CASE_ACCUSATIVE"
	CaseGenitive   caseClass = "CASE_GENITIVE"
	CaseJussive    caseClass = "CASE_JUSSIVE"
)
