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
	Id            string
	Word          string
	Letters       []LetterProps
	MoveRightUrl  string
	MoveLeftUrl   string
	AddWordUrl    string
	RemoveWordUrl string
}

type EditProps struct {
	ExcerptTitle string
	Username     string
	SelectedWord SelectedWordProps
	Words        []partials.WordProps
	Error        string
	Warning      string
	TitleUrl     string
	EditWordUrl  string
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
						Div(Class("flex justify-center"),
							partials.InspectorWordRegular(p.EditWordUrl, p.SelectedWord.Word),
						),
						Div(Class("flex flex-col gap-2 mx-2 "),
							Div(Class("flex justify-center gap-2"),
								FormEl(Class("w-full"), Method("post"), Action(p.SelectedWord.MoveRightUrl), up.Target("#text"),
									Button(Type("submit"), Class("w-full bg-sky-600 text-white rounded-lg p-2"),
										Img(Class("mx-auto h-5 invert"), Src("/static/icons/angles-right-solid.svg")),
									),
								),
								FormEl(Class("w-full"), Method("post"), Action(p.SelectedWord.MoveLeftUrl), up.Target("#text"),
									Button(Type("submit"), Class("w-full bg-sky-600 text-white rounded-lg p-2"),
										Img(Class("mx-auto h-5 invert"), Src("/static/icons/angles-left-solid.svg")),
									),
								),
							),
							FormEl(Class("w-full"), Method("post"), Action(p.SelectedWord.AddWordUrl),
								Button(Type("submit"), Class("w-full bg-sky-600 text-white rounded-lg p-2"),
									Img(Class("mx-auto h-5 invert"), Src("/static/icons/plus-solid.svg")),
								),
							),
							FormEl(Class("w-full"), Method("post"), Action(p.SelectedWord.RemoveWordUrl), g.Attr("onsubmit", `return confirm("Are you sure you want to delete this word?")`),
								Button(Type("submit"), Class("w-full bg-red-600 text-white rounded-lg p-2"),
									Img(Class("mx-auto h-5 invert"), Src("/static/icons/trash-solid.svg")),
								),
							),
						),
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
					partials.Text(p.Words),
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
