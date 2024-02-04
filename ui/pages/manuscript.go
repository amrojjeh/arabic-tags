package pages

import (
	"fmt"

	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type ManuscriptProps struct {
	ExcerptTitle        string
	ExcerptID           string
	Locked              bool
	AcceptedPunctuation string
	Content             string
	UnlockURL           string
	InspectorURL        string
	LockURL             string
}

func ManuscriptPage(p ManuscriptProps) g.Node {
	return MainBase(MainBaseProps{
		Title: fmt.Sprintf("%s - Arabic Tags", p.ExcerptTitle),
		Main: []g.Node{Div(Class("flex flex-col h-full"),
			g.If(p.Locked,
				Div(Class("text-center bg-yellow-300 text-black"),
					Img(Class("align-bottom"), Src("/static/warning.svg")),
					g.Text("This page is locked, so it cannot be modified. Click "),
					FormEl(Method("post"), Action(p.UnlockURL),
						Button(Class("underline cursor-pointer"), Type("submit"),
							g.Text("here"))),
					g.Text(" to unlock it (this resets all metadata)"))),
			H2(Class("text-2xl text-center"),
				g.Text(p.ExcerptTitle)),
			g.El("arabic-input", Class("grow"), ID(p.ExcerptID), g.Attr("punctuation", p.AcceptedPunctuation), Value(p.Content), g.If(p.Locked, ReadOnly())))},
		Nav: partials.MainNav(p.ExcerptID),
		Footer: []g.Node{
			Div(
				// TODO(Amr Ojjeh): Move to backend
				g.El("delete-errors"),
				g.El("delete-vowels"),
			),
			g.If(p.Locked,
				A(Href(p.InspectorURL),
					Button(Class("bg-sky-600 text-white rounded-lg p-2"),
						g.Text("Next"),
					)),
			),
			g.If(!p.Locked,
				FormEl(Method("post"), Action(p.LockURL),
					Button(Class("bg-sky-600 text-white rounded-lg p-2"),
						g.Text("Next"),
					),
				),
			),
		},
		Error: "",
	})
}
