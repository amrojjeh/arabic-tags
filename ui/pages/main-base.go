package pages

import (
	"github.com/amrojjeh/arabic-tags/ui"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type MainBaseProps struct {
	Title  string
	Main   []g.Node
	Nav    []g.Node
	Footer []g.Node
	Error  string
}

func MainBase(p MainBaseProps) g.Node {
	return ui.HTML5(ui.HTML5Props{
		Title:       p.Title,
		Description: "Tag arabic text with metadata and irab",
		Language:    "en",
		Head: []g.Node{
			Link(Rel("stylesheet"), Href("/static/main.css")),
			Script(Src("/static/main.js")),
		},
		Body: []g.Node{
			Class("h-screen gap-0 flex flex-col bg-red-50/25"),
			Nav(
				Class("text-white bg-red-800 px-5 py-2 grid grid-cols-3 grid-rows-1"),
				g.Group(p.Nav),
			),
			Div(
				Class("bg-red-200 text-center text-2xl text-red-800"),
				g.Text(p.Error),
			),
			// TODO(Amr Ojjeh): Make offline work
			Main(
				Class("grow p-0 overflow-y-hidden"),
				Div(ID("offline-warning"),
					Class("hidden text-center bg-yellow-300 text-black"),
					Img(Src("/static/warning.svg"), Class("inline align-bottom")),
					g.Text("You're offline. Any changes you make will not be saved until you're back online."),
				),
				g.Group(p.Main)),
			g.If(p.Footer != nil,
				Footer(
					Class("px-5 py-3 border-t-2 text-2xl justify-between flex"),
					g.Group(p.Footer),
				),
			),
		},
	})
}
