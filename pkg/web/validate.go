package web

import (
	"github.com/yfedoruck/todolist/pkg/env"
	"github.com/yfedoruck/todolist/pkg/lang"
	"github.com/yfedoruck/todolist/pkg/pg"
	"github.com/yfedoruck/todolist/pkg/resp"
	"html/template"
	"net/http"
	"path/filepath"
)

func isUniqueUsername(r *http.Request, data *RegisterData, db *pg.Postgres) bool {

	username := r.PostFormValue("username")

	ok := db.IsUniqueUsername(username)

	if ok {
		data.Error.Username = ""
	} else {
		data.Error.Username = lang.UsernameExists
	}

	return ok
}

func isUniqueEmail(r *http.Request, data *RegisterData, db *pg.Postgres) bool {

	ok := db.IsUniqueEmail(r.PostFormValue("email"))

	if ok {
		data.Error.Email = ""
	} else {
		data.Error.Email = lang.EmailExists
	}

	return ok
}

func clearRegisterForm(data *RegisterData) {
	data.Error = RegisterErr{}
	data.PreFill = RegisterField{}
}

func clearLoginForm(ld *LoginData) {
	ld.Error = ""
	ld.PreFill = LoginField{}
}

func renderTemplate(w http.ResponseWriter, tpl string, data interface{}) {
	viewDir := env.BasePath() + filepath.FromSlash("/views/")

	t, err := template.ParseFiles(viewDir + tpl + ".html")
	resp.Check(err)

	_ = t.Execute(w, data)
}
