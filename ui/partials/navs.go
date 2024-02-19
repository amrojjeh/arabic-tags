package partials

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func SimpleNav(username, loginUrl, registerUrl, logoutUrl string) []g.Node {
	return []g.Node{
		g.If(username == "", g.Group([]g.Node{H1(Class("font-black col-start-2 text-2xl text-center"),
			A(Href("/"),
				g.Text("Arabic Tags"),
			),
		),
			Div(Class("flex gap-3"),
				A(Class("underline text-center"), Href(loginUrl),
					g.Text("Login"),
				),
				A(Class("underline text-center"), Href(registerUrl),
					g.Text("Register"),
				),
			)})),
		g.If(username != "", g.Group([]g.Node{
			H1(Class("font-black col-start-2 text-2xl text-center"),
				A(Href("/"), g.Text("Arabic Tags")),
			),
			Div(Class("flex gap-3"),
				P(g.Text(username)),
				A(Class("underline"), g.Attr("up-method", "post"), Href(logoutUrl),
					g.Text("Logout"),
				),
			),
		})),
	}
}
