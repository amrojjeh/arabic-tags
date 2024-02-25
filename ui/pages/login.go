package pages

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type LoginProps struct {
	EmailField string
	EmailError string

	PasswordError string

	LoginError  string
	LoginUrl    string
	RegisterUrl string
	LogoutUrl   string
}

type LoginResponse struct {
	Email    string
	Password string
}

func LoginPage(p LoginProps) g.Node {
	return MainBase(MainBaseProps{
		Title: "Login",
		Nav:   partials.SimpleNav("", p.LoginError, p.RegisterUrl, p.LogoutUrl),
		Main: []g.Node{
			Div(Class("flex flex-col gap-3 items-center justify-center h-full"),
				H1(Class("text-2xl"),
					g.Text("Login"),
				),
				FormEl(Class("flex flex-col gap-3"), Action("/login"), Method("post"), AutoComplete("off"),
					g.If(p.LoginError != "", P(Class("w-60 text-md text-red-500 font-bold text-center pb-1"),
						g.Text(p.LoginError),
					)),
					g.If(p.EmailError != "", P(Class("w-60 text-md text-red-500 font-bold text-center"),
						g.Text(p.EmailError),
					)),
					Div(Class("flex gap-1"),
						Img(Class("w-4 opacity-30"), Src("/static/icons/envelope-solid.svg")),
						Input(Type("email"), Required(), Name("email"), Value(p.EmailField), Class("w-full block"), ID("email"), Placeholder("Enter email")),
					),
					g.If(p.PasswordError != "", P(Class("w-60 text-md text-red-500 font-bold text-center"),
						g.Text(p.PasswordError),
					)),
					Div(Class("flex gap-1"),
						Img(Class("w-4 opacity-30"), Src("/static/icons/lock-solid.svg")),
						Input(Required(), Name("password"), Type("password"), Class("w-full block"), ID("password"), Placeholder("Enter password")),
					),
					Button(Class("text-xl bg-sky-600 text-white rounded px-4 py-2"), Type("submit"),
						g.Text("Login"),
					),
				),
			),
		},
	})
}

func NewLoginResponse(r *http.Request) (LoginResponse, error) {
	err := r.ParseForm()
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}, nil
}

func (l LoginResponse) Props(loginUrl, registerUrl, logoutUrl string) LoginProps {
	return LoginProps{
		EmailField: l.Email,
	}
}
