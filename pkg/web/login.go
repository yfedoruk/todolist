package web

import (
	"errors"
	"github.com/yfedoruck/todolist/pkg/cookie"
	"github.com/yfedoruck/todolist/pkg/lang"
	"github.com/yfedoruck/todolist/pkg/pg"
	"github.com/yfedoruck/todolist/pkg/resp"
	"net/http"
)

type loginHandler struct {
	db   *pg.Postgres
	data *LoginData
}

type LoginData struct {
	Css     string
	Title   string
	Error   string
	PreFill LoginField
}

type LoginField struct {
	Username string
	Password string
}

func (l *loginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == nil {
		http.Redirect(w, r, "/todolist", http.StatusFound)
		return
	}

	if r.Method == "GET" {
		renderTemplate(w, "login", l.data)
	} else {
		err := r.ParseForm()
		resp.Check(err)
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")

		if username == "" {
			err = errors.New(lang.NameEmpty)
		}

		if password == "" {
			err = errors.New(lang.PassEmpty)
		}

		id, err := l.db.LoginUser(username, password)
		if err != nil {
			l.data.Error = err.Error()
			l.data.PreFill = LoginField{
				Username: username,
				Password: password,
			}
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			cookie.Cookie{
				Name: username,
				Id:   id,
			}.Set(w)
			clearLoginForm(l.data)
			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	}
}
