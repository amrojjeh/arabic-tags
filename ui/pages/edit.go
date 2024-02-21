package pages

import (
	"fmt"
	"strconv"

	"github.com/amrojjeh/arabic-tags/ui/partials"
	"github.com/amrojjeh/kalam"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type LetterProps struct {
	Letter     string
	ShortVowel rune
	Shadda     bool
	Index      int
	PostUrl    string
}

type SelectedWordProps struct {
	Word    string
	Letters []LetterProps
}

type WordProps struct {
	Word        string
	Punctuation bool
	Connected   bool
	Selected    bool
	GetUrl      string
}

type EditProps struct {
	ExcerptTitle string
	Username     string
	SelectedWord SelectedWordProps
	Words        []WordProps
	Error        string
	Warning      string
	TitleUrl     string
	ExportUrl    string
	LoginUrl     string
	RegisterUrl  string
	LogoutUrl    string
}

func EditPage(p EditProps) g.Node {
	return MainBase(MainBaseProps{
		Title: p.ExcerptTitle,
		Main: []g.Node{
			Div(Class("flex flex-col h-[99%]"),
				H2(Class("text-2xl flex justify-center"),
					partials.TitleRegular(p.TitleUrl, p.ExcerptTitle),
				),
				Div(g.Attr("dir", "rtl"), Class("grid grid-rows-1 grid-cols-[400px_auto] gap-4 h-[97%]"),
					Div(ID("inspector"), Class("border-e-2 m-1 h-full overflow-y-auto"),
						P(ID("inspector-word"), Class("text-5xl text-center leading-loose"),
							g.Text(p.SelectedWord.Word)),
						Div(Class("border-solid border-2 border-black bg-slate-200 align-center m-1 p-1"),
							P(
								g.Text("Stuff..."),
							),
						),
						g.Group(g.Map(p.SelectedWord.Letters, func(p LetterProps) g.Node {
							return FieldSet(c.Classes{
								"text-3xl m-1 p-4 leading-loose":      true,
								"border-dashed border-2 border-black": p.Index%2 != 0},
								Legend(Class("text-4xl"),
									g.Text(p.Letter),
								),
								FormEl(Method("post"), Action(p.PostUrl),
									Select(Name(strconv.Itoa(p.Index)), Class("block"), ID(fmt.Sprintf("letter-%v", p.Index)),
										Option(Value(string(kalam.Damma)), g.If(p.ShortVowel == kalam.Damma, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Damma)),
										),
										Option(Value(string(kalam.Dammatan)), g.If(p.ShortVowel == kalam.Dammatan, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Dammatan)),
										),
										Option(Value(string(kalam.Fatha)), g.If(p.ShortVowel == kalam.Fatha, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Fatha)),
										),
										Option(Value(string(kalam.Fathatan)), g.If(p.ShortVowel == kalam.Fathatan, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Fathatan)),
										),
										Option(Value(string(kalam.Kasra)), g.If(p.ShortVowel == kalam.Kasra, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Kasra)),
										),
										Option(Value(string(kalam.Kasratan)), g.If(p.ShortVowel == kalam.Kasratan, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Kasratan)),
										),
										Option(Value(string(kalam.Sukoon)), g.If(p.ShortVowel == kalam.Sukoon, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Sukoon)),
										),
									),
									Label(
										Input(Type("checkbox"), Value("true"), Name("shadda"), g.If(p.Shadda, Checked())),
										g.Text(string(kalam.Shadda)),
									),
								),
							)
						})),
					),
					Div(ID("text"),
						P(Class("text-4xl leading-loose"),
							g.Group(g.Map(p.Words, func(p WordProps) g.Node {
								if p.Punctuation {
									return Span(
										g.Text(p.Word),
										g.If(!p.Connected, g.Text(" ")),
									)
								}
								return A(Href(p.GetUrl), c.Classes{
									"cursor-pointer hover:text-red-700": !p.Selected,
									"text-sky-600":                      p.Selected,
								},
									g.Text(p.Word),
									g.If(!p.Connected, g.Text(" ")),
								)
							})),
						),
					),
				),
			),
		},
		Nav: partials.SimpleNav(p.Username, p.LoginUrl, p.RegisterUrl, p.LogoutUrl),
		Footer: []g.Node{
			A(Class("bg-sky-600 text-white rounded-lg p-2"), Href(p.ExportUrl), Target("_blank"),
				g.Text("Export"),
			),
		},
		Error:   p.Error,
		Warning: p.Warning,
	})
}
