package ui

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type HTML5Props struct {
	Title       string
	Description string
	Language    string
	HTMLClasses string
	Head        []g.Node
	Body        []g.Node
}

func HTML5(p HTML5Props) g.Node {
	return Doctype(
		HTML(g.If(p.Language != "", Lang(p.Language)), Class(p.HTMLClasses),
			Head(
				Meta(Charset("utf-8")),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1")),
				Link(Rel("preconnect"), Href("https://fonts.googleapis.com")),
				Link(Href("https://fonts.googleapis.com/css2?family=Noto+Sans+Arabic:wght@400..700&display=swap"), Rel("stylesheet")),
				Link(Rel("stylesheet"), Href("/static/unpoly.min.css")),
				Script(Src("/static/unpoly.min.js")),
				Link(Rel("stylesheet"), Href("/static/main.css")),
				Link(Rel("icon"), Type("image/x-icon"), Href("/static/icons/favicon.ico")),
				Script(Src("/static/main.js")),
				g.If(p.Description != "", Meta(Name("description"),
					Content(p.Description))),
				TitleEl(g.Text(p.Title)),
				g.Group(p.Head),
			),
			Body(g.Attr("up-main"), g.Group(p.Body)),
		),
	)
}

func SelectAttr() g.Node {
	return g.Attr("select")
}
