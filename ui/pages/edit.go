package pages

import (
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type EditProps struct {
	ExcerptTitle string
	Username     string
	Inspector    g.Node
	Text         g.Node
	Error        string
	Warning      string
	TitleUrl     string
	ExportUrl    string
	LoginUrl     string
	RegisterUrl  string
	LogoutUrl    string
}

func EditPage(p EditProps) g.Node {
	return MainBase(MainBaseProps{
		Title: p.ExcerptTitle,
		Main: []g.Node{
			Div(Class("flex flex-col h-[99%]"),
				H2(Class("text-2xl flex justify-center"),
					partials.TitleRegular(p.TitleUrl, p.ExcerptTitle),
				),
				Div(g.Attr("dir", "rtl"), Class("grid grid-rows-1 grid-cols-[400px_auto] gap-4 h-[97%]"),
					p.Inspector,
					p.Text,
				),
			),
		},
		Nav: partials.SimpleNav(p.Username, p.LoginUrl, p.RegisterUrl, p.LogoutUrl),
		Footer: []g.Node{
			A(Class("bg-sky-600 text-white rounded-lg p-2"), Href(p.ExportUrl), Target("_blank"),
				g.Text("Export"),
			),
		},
		Error:   p.Error,
		Warning: p.Warning,
	})
}
