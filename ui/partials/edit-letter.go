package partials

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	"github.com/amrojjeh/arabic-tags/ui"
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

func EditLetter(id, editUrl, selectUrl, selectedWord string, connected bool) g.Node {
	return Div(
		InspectorWordRegular(editUrl, selectedWord),
		TextWord(id, selectUrl, selectedWord, connected, true),
	)
}

func InspectorWordForm(editUrl, word string) g.Node {
	return FormEl(ID("inspector-word"), Action(editUrl), Method("post"), Class("flex items-center gap-2"),
		Input(AutoComplete("off"), ui.SelectAttr(), Name("new_word"), AutoFocus(), Value(word), Class("mr-2")),
		Button(Type("submit"), Img(Src("/static/icons/check-solid.svg"), Class("w-4"))),
	)
}

func InspectorWordRegular(editUrl, title string) g.Node {
	return A(ID("inspector-word"), Class("group flex items-center gap-2"), Href(editUrl), up.Target("#inspector-word"),
		P(Class("text-5xl text-center leading-loose"), g.Text(title)),
		Img(Src("/static/icons/pencil-solid.svg"), Class("inline w-4 invisible group-hover:visible")),
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
func Text(words []WordProps) g.Node {
	return Div(ID("text"),
		P(Class("text-4xl leading-loose"),
			g.Group(g.Map(words, func(p WordProps) g.Node {
				if p.Punctuation {
					return Span(
						g.Text(p.Word),
						g.If(!p.Connected, g.Text(" ")),
					)
				}
				return TextWord(p.Id, p.GetUrl, p.Word, p.Connected, p.Selected)
			})),
		),
	)
}

func KeyValueCheckbox(postUrl, key string, value bool) g.Node {
	return FormEl(Method("post"), Action(postUrl), up.AutoSubmit(), up.Target("#text"),
		P(Class("pe-2 text-2xl flex justify-end gap-2"),
			Input(Type("checkbox"), Name("value"), g.If(value, Checked())),
			g.Text(key),
		),
	)
}
