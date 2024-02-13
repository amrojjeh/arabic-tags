package pages

import (
	u "github.com/amrojjeh/arabic-tags/internal/unpoly"
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type HomeExcerpt struct {
	Name string
	Url  string
}

type HomeProps struct {
	Username string
	Excerpts []HomeExcerpt
	AddUrl   string
}

func HomePage(p HomeProps) g.Node {
	return MainBase(MainBaseProps{
		Title: "",
		Main: []g.Node{
			Div(Class("p-5 flex gap-5"),
				A(Class("flex flex-col items-center w-fit"), Href(p.AddUrl), u.Layer("new"), u.Mode("popup"), u.History(false),
					Div(Class("rounded-full bg-white h-7 w-7 border-dashed border-2 border-black flex items-center justify-center"),
						Img(Class("h-4 w-4"), Src("/static/icons/plus-solid.svg")),
					),
					P(g.Text("Add")),
				),
				g.Group(g.Map(p.Excerpts, func(e HomeExcerpt) g.Node {
					return A(Class("flex flex-col items-center w-fit"), Href(e.Url),
						Div(Class("rounded-full bg-red-500 h-7 w-7")),
						P(g.Text(e.Name)),
					)
				})),
			),
		},
		Nav: partials.UserNav(p.Username),
	})
}
