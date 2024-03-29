package pages

import (
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type ManuscriptProps struct {
	ExcerptTitle        string
	ReadOnly            bool
	AcceptedPunctuation string
	Content             string
	Warning             string
	Error               string
	Username            string
	SaveUrl             string
	NextUrl             string
	TitleUrl            string
	LoginUrl            string
	RegisterUrl         string
	LogoutUrl           string
}

func ManuscriptPage(p ManuscriptProps) g.Node {
	return MainBase(MainBaseProps{
		Title: p.ExcerptTitle,
		Main: []g.Node{
			Div(Class("flex flex-col h-full"),
				Div(Class("flex justify-center"),
					partials.TitleRegular(p.TitleUrl, p.ExcerptTitle, p.ReadOnly),
				),
				g.El("arabic-input", Class("grow"), g.Attr("url", p.SaveUrl), g.Attr("punctuation", p.AcceptedPunctuation), Value(p.Content), g.If(p.ReadOnly, ReadOnly())),
			)},
		Nav: partials.SimpleNav(p.Username, p.LoginUrl, p.RegisterUrl, p.LogoutUrl),
		Footer: []g.Node{
			Div(
				g.El("delete-errors"),
				g.El("delete-vowels"),
			),
			g.If(!p.ReadOnly,
				FormEl(Method("post"), Action(p.NextUrl), ID("next"),
					// TODO(Amr Ojjeh): Not secure yes, but sue me (after I fix it more properly later via web component)
					Button(Type("submit"), Class("bg-sky-600 text-white rounded-lg p-2"), g.Attr("onclick", "up.submit(`#next`); this.disabled = true"),
						g.Text("Next"),
					),
				),
			),
		},
		Error:   p.Error,
		Warning: p.Warning,
	})
}
