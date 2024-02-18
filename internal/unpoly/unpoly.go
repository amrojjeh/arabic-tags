package unpoly

import (
	g "github.com/maragudk/gomponents"
)

func Layer(mode string) g.Node {
	return g.Attr("up-layer", mode)
}

func Mode(mode string) g.Node {
	return g.Attr("up-mode", mode)
}

func History(mode bool) g.Node {
	if mode {
		return g.Attr("up-history", "true")
	}
	return g.Attr("up-history", "false")
}

func Target(selector string) g.Node {
	return g.Attr("up-target", selector)
}
