package pages

import (
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type HomeProps struct {
	TitleField string
	TitleError string

	PasswordError string
}

func HomePage(p HomeProps) g.Node {
	return MainBase(MainBaseProps{
		Title: "Home",
		Nav:   partials.SimpleNav(),
		Main: []g.Node{
			Div(Class("flex flex-col gap-3 items-center justify-center h-full"),
				H1(Class("text-2xl"),
					g.Text("Create a new excerpt"),
				),
				FormEl(Class("flex flex-col gap-3"), Action("/"), Method("post"), AutoComplete("off"),
					Div(
						P(Class("w-60 text-md text-red-500 font-bold text-center pb-1"),
							g.Text(p.TitleError),
						),
						Div(Class("flex gap-1"),
							Img(Class("w-4 opacity-30"), Src("/static/icons/heading-solid.svg")),
							Input(Required(), Name("title"), Value(p.TitleField), Class("w-full block"), ID("title"), Placeholder("Enter title")),
						),
					),
					Div(
						P(Class("w-60 text-md text-red-500 font-bold text-center pb-1"),
							g.Text(p.PasswordError),
						),
						Div(Class("flex gap-1"),
							Img(Class("w-4 opacity-30"), Src("/static/icons/lock-solid.svg")),
							Input(Required(), Name("password"), Type("password"), Class("w-full block"), ID("password"), Placeholder("Enter password")),
						),
					),
					Button(Class("text-xl bg-sky-600 text-white rounded px-4 py-2"), Type("submit"),
						g.Text("Create"),
					),
					P(Class("w-60"),
						g.Text("Once you create your excerpt, "),
						Span(Class("font-bold text-orange-600"),
							g.Text("BOOKMARK YOUR URL."),
						),
					),
					P(Class("w-60"),
						g.Text("There's no other way to access your excerpt unless you have the URL."),
					),
				),
			),
		},
	})
}
