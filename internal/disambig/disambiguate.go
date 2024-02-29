package disambig

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

type Word struct {
	Word        string
	Connected   bool
	Punctuation bool
}

// TODO(Amr Ojjeh): Use goroutines
func Disambiguate(text string) ([]Word, error) {
	texts := splitText(text)
	var words []Word
	for _, t := range texts {
		res, err := request(t)
		if err != nil {
			return nil, err
		}

		words = make([]Word, 0, len(res.Output.Disambig))
		for _, cWord := range res.Output.Disambig {
			a := cWord.Analyses[0].Analysis
			if a.Pos == "punc" {
				words = append(words, Word{
					Word:        a.ATBTok,
					Connected:   false,
					Punctuation: true,
				})
			} else {
				ts := strings.Split(strings.ReplaceAll(a.ATBTok, "_", ""), "+")
				ts_len := len(ts)
				for i, t := range ts {
					words = append(words, Word{
						Word:        t,
						Connected:   ts_len > 1 && i != ts_len-1,
						Punctuation: false,
					})
				}
			}
		}
	}
	return words, nil
}

func request(text string) (camelResponse, error) {
	data, err := json.Marshal(camelRequest{
		Dialect:  "msa",
		Sentence: text,
	})
	if err != nil {
		return camelResponse{}, err
	}

	body := bytes.NewReader(data)
	resp, err := http.Post("https://camelira.abudhabi.nyu.edu/api/disambig",
		"application/json", body)
	if err != nil {
		return camelResponse{}, errors.Join(ErrRequest, err)
	}
	res, err := io.ReadAll(resp.Body)
	if err != nil {
		return camelResponse{}, errors.Join(ErrBadResponse, err)
	}

	inter := camelResponse{}
	err = json.Unmarshal(res, &inter)
	if err != nil {
		return camelResponse{}, errors.Join(BadFormatError{Text: string(res),
			ExpectedFormat: "json"}, err)
	}

	return inter, nil
}

func splitText(text string) []string {
	texts := []string{}
	for utf8.RuneCountInString(text) > 400 {
		words := strings.Split(text, " ")
		running := 0
		lastIndex := 0 // exclusive
		for running < 400 {
			// Adding one to account for space
			running += utf8.RuneCountInString(words[lastIndex]) + 1
			lastIndex += 1
		}
		texts = append(texts, strings.Join(words[:lastIndex], " "))
		text = strings.Join(words[lastIndex:], " ")
	}

	if text != "" {
		texts = append(texts, text)
	}

	return texts
}

type camelRequest struct {
	Dialect string `json:"dialect"`
	Sentence string `json:"sentence"`
}

type camelResponse struct {
	Output camelOutput `json:"output"`
}

type camelOutput struct {
	Disambig []camelWord `json:"disambig"`
}

type camelWord struct {
	Analyses []camelAnalysisMeta `json:"analyses"`
}

type camelAnalysisMeta struct {
	Analysis camelAnalysis `json:"analysis"`
}

type camelAnalysis struct {
	Pos    string `json:"pos"`
	ATBTok string `json:"atbtok"`
}
