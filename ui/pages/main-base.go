package pages

import (
	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type MainBaseProps struct {
	Title   string
	Main    []g.Node
	Nav     []g.Node
	Footer  []g.Node
	Error   string
	Warning string
}

func MainBase(p MainBaseProps) g.Node {
	return ui.HTML5(ui.HTML5Props{
		Title:       p.Title,
		Description: "Tag arabic text with metadata and irab",
		Language:    "en",
		Body: []g.Node{
			Class("h-screen gap-0 flex flex-col bg-red-50/25"),
			g.If(p.Nav != nil, Nav(Class("text-white bg-red-800 px-5 py-2 grid grid-cols-3 grid-rows-1 items-center"),
				g.Group(p.Nav),
			)),
			partials.Error(p.Error),
			Main(Class("grow p-0 overflow-y-hidden"),
				g.If(p.Warning != "", Div(ID("any-warning"),
					Class("text-center bg-yellow-300 text-black"),
					Img(Src("/static/icons/warning.svg"), Class("inline align-bottom")),
					g.Text(p.Warning),
				)),
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
