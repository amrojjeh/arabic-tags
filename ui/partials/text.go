package partials

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type WordProps struct {
	Id          string
	Word        string
	Punctuation bool
	Connected   bool
	Selected    bool
	GetUrl      string
}

func Text(words []WordProps) g.Node {
	return Div(ID("text"),
		P(Class("text-4xl leading-loose"),
			g.Group(g.Map(words, func(p WordProps) g.Node {
				return TextWord(p.Id, p.GetUrl, p.Word, p.Connected, p.Selected)
			})),
		),
	)
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
