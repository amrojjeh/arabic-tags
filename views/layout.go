package views

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type Prop struct {
	Title   string
	Content g.Node
}

func Page(p Prop) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    p.Title,
		Language: "en",
		Head:     []g.Node{},
		Body: []g.Node{
			P(g.Text("This is a test")),
		},
	})
}
