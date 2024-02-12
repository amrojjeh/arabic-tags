package partials

import (
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func SimpleNav() []g.Node {
	return []g.Node{
		Div(),
		H1(Class("font-black col-start-2 text-2xl text-center"),
			A(Href("/"),
				g.Text("Arabic Tags"),
			),
		),
		Div(Class("flex gap-3"),
			A(Class("underline text-center"), Href("/login"),
				g.Text("Login"),
			),
			A(Class("underline text-center"), Href("/register"),
				g.Text("Register"),
			),
		),
	}
}

func MainNav(ID string) []g.Node {
	return []g.Node{
		H2(Class("select-none col-start-1 text-xl text-center"),
			g.Text("saving..."),
		),
		H1(Class("font-black col-start-2 text-2xl text-center"),
			A(Href("/"), g.Text("Arabic Tags")),
		),
		Div(Class("flex gap-0"),
			Div(Class("flex items-center bg-sky-600 text-md ps-1 pe-1 h-full"),
				P(g.Text("ID")),
			),
			// TODO(Amr Ojjeh): Write javascript in main.js
			// onclick="this.select();navigator.clipboard.writeText('{{.}}')"
			Input(Class("text-xs text-black cursor-pointer"), Type("text"), Value(ID), ReadOnly()),
		),
	}
}

func UserNav(username string) []g.Node {
	return []g.Node{
		Div(),
		H1(Class("font-black col-start-2 text-2xl text-center"),
			A(Href("/"), g.Text("Arabic Tags")),
		),
		Div(Class("flex gap-3"),
			P(g.Text(username)),
			A(Class("underline"), g.Attr("up-method", "post"), Href("/logout"),
				g.Text("Logout"),
			),
		),
	}
}
