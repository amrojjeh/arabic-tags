package pages

import (
	"net/http"

	"github.com/amrojjeh/arabic-tags/ui/partials"
	g "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

type RegisterProps struct {
	UsernameField string
	UsernameError string

	EmailField string
	EmailError string

	PasswordError string
}

type RegisterResponse struct {
	Username   string
	Email      string
	Password   string
	RePassword string
}

func RegisterPage(p RegisterProps) g.Node {
	return MainBase(MainBaseProps{
		Title: "Register",
		Nav:   partials.SimpleNav(),
		Main: []g.Node{
			Div(Class("flex flex-col gap-3 items-center justify-center h-full"),
				H1(Class("text-2xl"),
					g.Text("Register an account"),
				),
				FormEl(Class("flex flex-col gap-3"), Action("/register"), Method("post"), AutoComplete("off"),
					g.If(p.UsernameError != "", P(Class("w-60 text-md text-red-500 font-bold text-center"),
						g.Text(p.UsernameError),
					)),
					Div(Class("flex gap-1"),
						Img(Class("w-4 opacity-30"), Src("/static/icons/user-solid.svg")),
						Input(Required(), Name("username"), Value(p.UsernameField), Class("w-full block"), ID("username"), Placeholder("Enter username")),
					),
					g.If(p.EmailError != "", P(Class("w-60 text-md text-red-500 font-bold text-center"),
						g.Text(p.EmailError),
					)),
					Div(Class("flex gap-1"),
						Img(Class("w-4 opacity-30"), Src("/static/icons/envelope-solid.svg")),
						Input(Required(), Name("email"), Value(p.EmailField), Class("w-full block"), ID("email"), Placeholder("Enter email")),
					),
					g.If(p.PasswordError != "", P(Class("w-60 text-md text-red-500 font-bold text-center"),
						g.Text(p.PasswordError),
					)),
					Div(Class("flex gap-1"),
						Img(Class("w-4 opacity-30"), Src("/static/icons/lock-solid.svg")),
						Input(Required(), Name("password"), Type("password"), Class("w-full block"), ID("password"), Placeholder("Enter password")),
					),
					Div(Class("flex gap-1"),
						Img(Class("w-4 opacity-30"), Src("/static/icons/lock-solid.svg")),
						Input(Required(), Name("repassword"), Type("password"), Class("w-full block"), ID("repassword"), Placeholder("Renter password")),
					),
					Button(Class("text-xl bg-sky-600 text-white rounded px-4 py-2"), Type("submit"),
						g.Text("Register"),
					),
				),
			),
		},
	})
}

func NewRegisterResponse(r *http.Request) (RegisterResponse, error) {
	err := r.ParseForm()
	if err != nil {
		return RegisterResponse{}, err
	}

	return RegisterResponse{
		Username:   r.Form.Get("username"),
		Email:      r.Form.Get("email"),
		Password:   r.Form.Get("password"),
		RePassword: r.Form.Get("repassword"),
	}, nil
}

func (r RegisterResponse) Props() RegisterProps {
	return RegisterProps{
		UsernameField: r.Username,
		EmailField:    r.Email,
	}
}
