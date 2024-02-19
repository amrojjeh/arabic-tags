package partials

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	"github.com/amrojjeh/arabic-tags/ui"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func TitleForm(postUrl, title string) g.Node {
	return FormEl(ID("excerpt-title"), Action(postUrl), Method("post"), up.Target("#excerpt-title"), Class("flex items-center"),
		Input(AutoComplete("off"), ui.SelectAttr(), Name("title"), AutoFocus(), Value(title), Class("mr-2")),
		Button(Type("submit"), Img(Src("/static/icons/check-solid.svg"), Class("w-4"))),
	)
}

func TitleRegular(getUrl, title string) g.Node {
	return A(ID("excerpt-title"), Class("group flex items-center"), Href(getUrl), up.Target("#excerpt-title"),
		P(Class("text-2xl pr-2 inline"), g.Text(title)),
		Img(Src("/static/icons/pencil-solid.svg"), Class("inline w-4 invisible group-hover:visible")),
	)
}
