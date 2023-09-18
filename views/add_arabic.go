package views

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func AddArabic() g.Node {
	return Main(Class("grow"),
		g.El("arabic-text"))
}
