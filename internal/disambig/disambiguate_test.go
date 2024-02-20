package disambig

import (
	"fmt"
	"testing"

	"github.com/amrojjeh/kalam"
	"github.com/amrojjeh/kalam/assert"
)

func TestDisambiguate(t *testing.T) {
	words, err := Disambiguate(kalam.FromBuckwalter("h*A byth."))
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(words), 4)

	assert.Equal(t, words[0].Word, fmt.Sprint(string(kalam.Heh), string(kalam.SuperscriptAlef), string(kalam.Thal), string(kalam.Alef)))
	assert.Equal(t, words[0].Connected, false)
	assert.Equal(t, words[0].Punctuation, false)

	assert.Equal(t, words[1].Word, kalam.FromBuckwalter("bayotu"))
	assert.Equal(t, words[1].Connected, true)
	assert.Equal(t, words[1].Punctuation, false)

	assert.Equal(t, words[2].Word, kalam.FromBuckwalter("hu"))
	assert.Equal(t, words[2].Connected, false)
	assert.Equal(t, words[2].Punctuation, false)

	assert.Equal(t, words[3].Word, ".")
	assert.Equal(t, words[3].Connected, false)
	assert.Equal(t, words[3].Punctuation, true)
}
