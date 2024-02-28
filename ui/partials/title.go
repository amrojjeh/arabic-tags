package partials

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	"github.com/amrojjeh/arabic-tags/ui"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func TitleForm(cancelUrl, postUrl, title string) g.Node {
	return FormEl(ID("excerpt-title"), ui.SelectAttr(), Action(postUrl), Method("post"), up.Target("#excerpt-title"), Class("flex items-center gap-2"),
		Input(AutoComplete("off"), Name("title"), AutoFocus(), Value(title), Class("mr-2")),
		A(Href(cancelUrl), Img(Src("/static/icons/xmark-solid.svg"), Class("w-4"))),
		Button(Type("submit"), Img(Src("/static/icons/check-solid.svg"), Class("w-4"))),
	)
}

func TitleRegular(getUrl, title string, readonly bool) g.Node {
	if getUrl != "" && !readonly {
		return A(ID("excerpt-title"), Class("group flex items-center"), Href(getUrl), up.Target("#excerpt-title"),
			P(Class("text-2xl pr-2 inline"), g.Text(title)),
			Img(Src("/static/icons/pencil-solid.svg"), Class("inline w-4 invisible group-hover:visible")),
		)
	}
	return A(ID("excerpt-title"), Class("group flex items-center"),
		P(Class("text-2xl pr-2 inline"), g.Text(title)),
	)
}
