package pages

import (
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type HomeProps struct {
	Username string
}

func HomePage(p HomeProps) g.Node {
	return MainBase(MainBaseProps{
		Title: "",
		Main: []g.Node{
			Div(g.Text("Home Page")),
		},
		Nav: partials.UserNav(p.Username),
	})
}
