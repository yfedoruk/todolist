package web

import (
	"github.com/yfedoruck/todolist/pkg/cookie"
	"github.com/yfedoruck/todolist/pkg/env"
	"github.com/yfedoruck/todolist/pkg/pg"
	"github.com/yfedoruck/todolist/pkg/resp"
	"github.com/yfedoruck/todolist/pkg/validate"
	"html/template"
	"net/http"
	"path/filepath"
)

type registerHandler struct {
	db   *pg.Postgres
	data *RegisterData
}

type RegisterData struct {
	Css     string
	Title   string
	Error   RegisterErr
	PreFill RegisterField
}
type RegisterErr struct {
	Email    string
	Username string
	Password string
}
type RegisterField struct {
	Email    string
	Username string
	Password string
}

func (h *registerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		viewDir := env.BasePath() + filepath.FromSlash("/views/")
		t, err := template.ParseFiles(viewDir + "register.html")
		resp.Check(err)

		err = t.Execute(w, h.data)
		resp.Check(err)
	} else {
		var validation = true

		err := r.ParseForm()
		resp.Check(err)
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		email := r.PostFormValue("email")

		errUsername := validate.Username(r)
		if errUsername != nil {
			h.data.Error.Username = errUsername.Error()
			validation = false
		} else {
			h.data.Error.Username = ""
		}

		errEmail := validate.Email(r)
		if errEmail != nil {
			h.data.Error.Email = errEmail.Error()
			validation = false
		} else {
			h.data.Error.Email = ""
		}

		if !validation {
			h.data.PreFill = RegisterField{
				email,
				username,
				password,
			}

			http.Redirect(w, r, "/register", http.StatusFound)
			return
		}

		isUniqueUsername := isUniqueUsername(r, h.data, h.db)
		isUniqueEmail := isUniqueEmail(r, h.data, h.db)

		if !isUniqueUsername || !isUniqueEmail {
			h.data.PreFill = RegisterField{
				email,
				username,
				password,
			}
			http.Redirect(w, r, "/register", http.StatusFound)
			return
		} else {
			h.data.PreFill = RegisterField{}
			h.data.Error = RegisterErr{}
		}

		id := h.db.RegisterUser(password, username, email)

		cookie.Cookie{
			Name: username,
			Id:   id,
		}.Set(w)

		clearRegisterForm(h.data)
		http.Redirect(w, r, "/todolist", http.StatusFound)
	}
}
