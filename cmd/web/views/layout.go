package views

import (
	g "github.com/maragudk/gomponents"
	c "github.com/maragudk/gomponents/components"
	. "github.com/maragudk/gomponents/html"
)

type Prop struct {
	Title   string
	Content g.Node
}

func Page(p Prop) g.Node {
	return c.HTML5(c.HTML5Props{
		Title:    p.Title,
		Language: "en",
		Head: []g.Node{
			Link(Type("text/css"), Rel("stylesheet"), Href("static/main.css")),
			Link(Rel("preconnect"), Href("https://fonts.googleapis.com")),
			Link(Rel("preconnect"), Href("https://fonts.gstatic.com"),
				g.Attr("crossorigin")),
			Link(Rel("stylesheet"),
				Href("https://fonts.googleapis.com/css2?family=Amiri&display=swap")),
			Script(Type("module"), Src("static/main.js")),
		},
		Body: []g.Node{Class("h-screen flex flex-col p4 bg-red-50/25"),
			navBar(p),
			p.Content,
		},
	})
}

func navBar(p Prop) g.Node {
	return Nav(Class(
		"text-white bg-red-800 px-5 py-2 flex justify-between items-center text-2xl "+
			"border-b-2 drop-shadow-md"),
		A(Href("#"), Class("underline"), g.Text("Back")),
		P(Class("font-black"), g.Text("Arabic Text")),
		Button(Type("button"),
			Class("bg-sky-600 px-3 py-2 rounded-lg"),
			g.Text("Share")))
}
