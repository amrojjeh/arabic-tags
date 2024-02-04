package pages

import (
	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type HomeProps struct {
	TitleField string
	TitleError string
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
				FormEl(Class("flex flex-col gap-3"), Action("/"), Method("post"),
					Div(
						P(Class("text-md text-red-500 font-bold text-center pb-1"),
							g.Text(p.TitleError),
						),
						Input(Required(), Name("title"), Value(p.TitleField), Class("block"), ID("title"), Placeholder("Enter title")),
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
