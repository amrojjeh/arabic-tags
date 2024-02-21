package partials

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

func EditLetter(id, getUrl, selectedWord string, connected bool) g.Node {
	return Div(
		InspectorWord(selectedWord),
		TextWord(id, getUrl, selectedWord, connected, true),
	)
}

func InspectorWord(word string) g.Node {
	return P(ID("inspector-word"), Class("text-5xl text-center leading-loose"),
		g.Text(word))
}

func TextWord(id, getUrl, word string, connected, selected bool) g.Node {
	return A(ID("w"+id), Href(getUrl), up.History(false), c.Classes{
		"cursor-pointer hover:text-red-700": !selected,
		"text-sky-600":                      selected,
	},
		g.Text(word),
		g.If(!connected, g.Text(" ")),
	)
}
