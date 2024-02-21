package partials

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func Error(err string) g.Node {
	return Div(ID("main-error"), Class("bg-red-200 text-center text-2xl text-red-800"), up.Hungry(),
		g.Text(err),
	)
}

func WithError(err string, node ...g.Node) g.Node {
	return Div(
		Error(err),
		g.Group(node),
	)
}
