package pages

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	"github.com/amrojjeh/arabic-tags/ui/partials"
	"github.com/amrojjeh/kalam"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type LetterProps struct {
	Letter          string
	ShortVowel      rune
	Shadda          bool
	SuperscriptAlef bool
	Index           int
	PostUrl         string
}

type SelectedWordProps struct {
	Id      string
	Word    string
	Letters []LetterProps
}

type WordProps struct {
	Id          string
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
						partials.InspectorWord(p.SelectedWord.Word),
						Div(Class("border-solid border-2 border-black bg-slate-200 align-center m-1 p-1"),
							P(
								g.Text("Stuff..."),
							),
						),
						g.Group(g.Map(p.SelectedWord.Letters, func(lp LetterProps) g.Node {
							return FieldSet(c.Classes{
								"text-3xl m-1 p-4 leading-loose":      true,
								"border-dashed border-2 border-black": lp.Index%2 != 0},
								Legend(Class("text-4xl"),
									g.Text(lp.Letter),
								),
								FormEl(Method("post"), Action(lp.PostUrl), up.AutoSubmit(), up.Target("#inspector-word, #w"+p.SelectedWord.Id),
									Select(Name("vowel"), Class("block"),
										Option(Value(string(kalam.Damma)), g.If(lp.ShortVowel == kalam.Damma, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Damma)),
										),
										Option(Value(string(kalam.Dammatan)), g.If(lp.ShortVowel == kalam.Dammatan, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Dammatan)),
										),
										Option(Value(string(kalam.Fatha)), g.If(lp.ShortVowel == kalam.Fatha, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Fatha)),
										),
										Option(Value(string(kalam.Fathatan)), g.If(lp.ShortVowel == kalam.Fathatan, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Fathatan)),
										),
										Option(Value(string(kalam.Kasra)), g.If(lp.ShortVowel == kalam.Kasra, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Kasra)),
										),
										Option(Value(string(kalam.Kasratan)), g.If(lp.ShortVowel == kalam.Kasratan, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Kasratan)),
										),
										Option(Value(string(kalam.Sukoon)), g.If(lp.ShortVowel == kalam.Sukoon, Selected()),
											g.Text(string(kalam.Placeholder)+string(kalam.Sukoon)),
										),
									),
									Div(Class("flex gap-2"),
										Label(
											Input(Type("checkbox"), Value("true"), Name("shadda"), g.If(lp.Shadda, Checked())),
											g.Text(string(kalam.Shadda)),
										),
										Label(
											Input(Type("checkbox"), Value("true"), Name("superscript_alef"), g.If(lp.SuperscriptAlef, Checked())),
											g.Text(string(kalam.SuperscriptAlef)),
										),
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
								return partials.TextWord(p.Id, p.GetUrl, p.Word, p.Connected, p.Selected)
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
