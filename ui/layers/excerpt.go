package layers

import (
	"net/http"

	u "github.com/amrojjeh/arabic-tags/internal/unpoly"

	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type ExcerptResponse struct {
	Title string
}

func ExcerptLayer(postUrl string) g.Node {
	return LayerBase(LayerBaseProps{
		Title: "Create excerpt",
		Main: []g.Node{
			FormEl(Class("flex flex-col gap-2"), Method("post"), Action(postUrl), u.Layer("root"),
				H1(Class("text-center text-xl"), g.Text("Create excerpt")),
				Input(AutoComplete("off"), Class("border-2 border-solid"), Name("title"), AutoFocus(), Type("text"), Required(), Placeholder("Enter title")),
				Button(Class("text-xl bg-sky-600 text-white rounded px-4 py-2"), Type("submit"),
					g.Text("Create"),
				),
			),
		},
	})
}

func NewExcerptResponse(r *http.Request) (ExcerptResponse, error) {
	err := r.ParseForm()
	if err != nil {
		return ExcerptResponse{}, err
	}
	return ExcerptResponse{
		Title: r.Form.Get("title"),
	}, nil
}
