package pages

import (
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type EditProps struct {
	ExcerptTitle string
	Inspector    g.Node
	Text         g.Node
	Nav          []g.Node
	ReadOnly     bool
	Error        string
	Warning      string
	TitleUrl     string
	ExportUrl    string
}

func EditPage(p EditProps) g.Node {
	return MainBase(MainBaseProps{
		Title: p.ExcerptTitle,
		Main: []g.Node{
			Div(Class("flex flex-col h-[99%]"),
				H2(Class("text-2xl flex justify-center"),
					partials.TitleRegular(p.TitleUrl, p.ExcerptTitle, p.ReadOnly),
				),
				Div(g.Attr("dir", "rtl"), Class("grid grid-rows-1 grid-cols-[400px_auto] gap-4 h-[97%]"),
					p.Inspector,
					p.Text,
				),
			),
		},
		Nav: p.Nav,
		Footer: []g.Node{
			A(Class("bg-sky-600 text-white rounded-lg p-2"), Href(p.ExportUrl), Target("_blank"),
				g.Text("Export"),
			),
		},
		Error:   p.Error,
		Warning: p.Warning,
	})
}
