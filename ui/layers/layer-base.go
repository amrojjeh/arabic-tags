package layers

import (
	"github.com/amrojjeh/arabic-tags/ui"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type LayerBaseProps struct {
	Title string
	Main  []g.Node
}

func LayerBase(p LayerBaseProps) g.Node {
	return ui.HTML5(ui.HTML5Props{
		Title:       p.Title,
		Description: "Tag arabic text with metadata and irab",
		Language:    "en",
		Body: []g.Node{
			Class("gap-0 flex flex-col"),
			Main(Class("grow p-4 overflow-y-hidden"),
				g.Group(p.Main),
			),
		},
	})
}
