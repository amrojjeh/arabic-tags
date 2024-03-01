package partials

import (
	up "github.com/amrojjeh/arabic-tags/internal/unpoly"
	"github.com/amrojjeh/arabic-tags/ui"
	"github.com/amrojjeh/kalam"
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type InspectorProps struct {
	Id               string
	Word             string
	Letters          []LetterProps
	Connected        bool
	Ignore           bool
	SentenceStart    bool
	CaseOptions      []DropdownOption
	StateOptions     []DropdownOption
	ReadOnly         bool
	EditWordUrl      string
	CaseUrl          string
	StateUrl         string
	MoveRightUrl     string
	MoveLeftUrl      string
	AddWordUrl       string
	RemoveWordUrl    string
	ConnectedUrl     string
	IgnoreUrl        string
	SentenceStartUrl string
}

type LetterProps struct {
	Letter          string
	ShortVowel      rune
	Shadda          bool
	SuperscriptAlef bool
	Index           int
	PostUrl         string
}

func EditLetter(id, editUrl, selectUrl, selectedWord string, connected bool) g.Node {
	return Div(
		InspectorWordRegular(editUrl, selectedWord, false),
		TextWord(id, selectUrl, selectedWord, connected, true),
	)
}

func Inspector(p InspectorProps) g.Node {
	return Div(ID("inspector"), Class("border-e-2 m-1 h-full overflow-y-auto"),
		Div(Class("flex justify-center"),
			InspectorWordRegular(p.EditWordUrl, p.Word, p.ReadOnly),
		),
		g.If(!p.ReadOnly, Div(Class("flex flex-col gap-2 mx-2 "),
			Div(Class("flex justify-center gap-2"),
				FormEl(Class("w-full"), Method("post"), Action(p.MoveRightUrl), up.Target("#text"),
					Button(Type("submit"), Class("w-full bg-sky-600 text-white rounded-lg p-2"),
						Img(Class("mx-auto h-5 invert"), Src("/static/icons/angles-right-solid.svg")),
					),
				),
				FormEl(Class("w-full"), Method("post"), Action(p.MoveLeftUrl), up.Target("#text"),
					Button(Type("submit"), Class("w-full bg-sky-600 text-white rounded-lg p-2"),
						Img(Class("mx-auto h-5 invert"), Src("/static/icons/angles-left-solid.svg")),
					),
				),
			),
			FormEl(Class("w-full"), Method("post"), Action(p.AddWordUrl),
				Button(Type("submit"), Class("w-full bg-sky-600 text-white rounded-lg p-2"),
					Img(Class("mx-auto h-5 invert"), Src("/static/icons/plus-solid.svg")),
				),
			),
			FormEl(Class("w-full"), Method("post"), Action(p.RemoveWordUrl), g.Attr("onsubmit", `return confirm("Are you sure you want to delete this word?")`),
				Button(Type("submit"), Class("w-full bg-red-600 text-white rounded-lg p-2"),
					Img(Class("mx-auto h-5 invert"), Src("/static/icons/trash-solid.svg")),
				),
			),
		)),
		Div(Class("border-solid border-2 border-black bg-slate-200 align-center m-2 p-1"),
			g.If(p.ConnectedUrl != "", KeyValueCheckbox(p.ConnectedUrl, "Connected", p.Connected, p.ReadOnly)),
			g.If(p.SentenceStartUrl != "", KeyValueCheckbox(p.SentenceStartUrl, "Sentence Start", p.SentenceStart, p.ReadOnly)),
			g.If(p.IgnoreUrl != "", KeyValueCheckbox(p.IgnoreUrl, "Ignore", p.Ignore, p.ReadOnly)),

			Div(Class("flex flex-col py-1 gap-1"),
				g.If(len(p.CaseOptions) != 0, KeyValueDropdown(p.CaseUrl, "Case", p.CaseOptions, p.ReadOnly)),
				g.If(len(p.StateOptions) != 0, KeyValueDropdown(p.StateUrl, "State", p.StateOptions, p.ReadOnly)),
			),
		),
		g.If(!p.ReadOnly, g.Group(g.Map(p.Letters, func(lp LetterProps) g.Node {
			return FieldSet(c.Classes{
				"text-3xl m-1 p-4 leading-loose":      true,
				"border-dashed border-2 border-black": lp.Index%2 != 0},
				Legend(Class("text-4xl"),
					g.Text(lp.Letter),
				),
				FormEl(Method("post"), Action(lp.PostUrl), up.AutoSubmit(), up.Target("#inspector-word, #w"+p.Id),
					Select(Name("vowel"), Class("block"),
						Option(Value("blank"), g.If(lp.ShortVowel == 0, Selected()),
							g.Text(""),
						),
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
		}))),
	)

}

func InspectorWordForm(cancelUrl, editUrl, id, word string) g.Node {
	return FormEl(ID("inspector-word"), ui.SelectAttr(), Action(editUrl), Method("post"), Class("flex items-center gap-2"),
		Input(AutoComplete("off"), Name("new_word"), AutoFocus(), Value(word), Class("mr-2")),
		A(Href(cancelUrl), Img(Src("/static/icons/xmark-solid.svg"), Class("w-4"))),
		Button(Type("submit"), Img(Src("/static/icons/check-solid.svg"), Class("w-4"))),
	)
}

func InspectorWordRegular(editUrl, title string, readonly bool) g.Node {
	if editUrl != "" && !readonly {
		return A(ID("inspector-word"), Class("group flex items-center gap-2"), Href(editUrl), up.Target("#inspector-word"),
			P(Class("text-5xl text-center leading-loose"), g.Text(title)),
			Img(Src("/static/icons/pencil-solid.svg"), Class("inline w-4 invisible group-hover:visible")),
		)
	}
	return A(ID("inspector-word"), Class("group flex items-center gap-2"),
		P(Class("text-5xl text-center leading-loose"), g.Text(title)),
	)
}

func KeyValueCheckbox(postUrl, key string, value, readonly bool) g.Node {
	if !readonly {
		return FormEl(Method("post"), Action(postUrl), up.AutoSubmit(), up.Target("#text"),
			P(Class("pe-2 text-2xl flex justify-end gap-2"),
				g.Text(key),
				Input(Type("checkbox"), Name("value"), g.If(value, Checked())),
			),
		)
	}
	return P(Class("pe-2 text-2xl flex justify-end gap-2"),
		g.Text(key),
		Input(Type("checkbox"), g.If(value, Checked()), Disabled()),
	)
}

type DropdownOption struct {
	Value    string
	Selected bool
}

func KeyValueDropdown(postUrl, key string, values []DropdownOption, readonly bool) g.Node {
	if !readonly {
		return FormEl(Method("post"), Action(postUrl), up.AutoSubmit(), up.Target("#inspector"),
			P(Class("pe-2 text-2xl flex justify-end gap-2"),
				Select(Name("value"), Class("bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block px-1"),
					g.Group(g.Map(values, func(d DropdownOption) g.Node {
						return Option(Value(d.Value), g.Text(d.Value), g.If(d.Selected, Selected()))
					})),
				),
				g.Text(key),
			),
		)
	}

	return P(Class("pe-2 text-2xl flex justify-end gap-2"),
		Select(Class("bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block px-1"), Disabled(),
			g.Group(g.Map(values, func(d DropdownOption) g.Node {
				return Option(Value(d.Value), g.Text(d.Value), g.If(d.Selected, Selected()))
			})),
		),
		g.Text(key),
	)
}
