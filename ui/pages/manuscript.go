package pages

import (
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type ManuscriptProps struct {
	ExcerptTitle        string
	ReadOnly            bool
	AcceptedPunctuation string
	Content             string
	SubmitUrl           string
	Warning             string
	Error               string
	Username            string
}

func ManuscriptPage(p ManuscriptProps) g.Node {
	return MainBase(MainBaseProps{
		Title: p.ExcerptTitle,
		Main: []g.Node{Div(Class("flex flex-col h-full"),
			H2(Class("text-2xl text-center"),
				g.Text(p.ExcerptTitle)),
			g.El("arabic-input", Class("grow"), g.Attr("url", p.SubmitUrl), g.Attr("punctuation", p.AcceptedPunctuation), Value(p.Content), g.If(p.ReadOnly, ReadOnly())))},
		Nav: partials.SimpleNav(p.Username),
		Footer: []g.Node{
			Div(
				// TODO(Amr Ojjeh): Move to backend
				g.El("delete-errors"),
				g.El("delete-vowels"),
			),
			g.If(!p.ReadOnly,
				FormEl(Method("post"), Action(p.SubmitUrl),
					Button(Class("bg-sky-600 text-white rounded-lg p-2"),
						g.Text("Next"),
					),
				),
			),
		},
		Error:   p.Error,
		Warning: p.Warning,
	})
}
