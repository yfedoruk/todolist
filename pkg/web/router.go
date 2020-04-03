package web

import (
	"github.com/yfedoruck/todolist/pkg/env"
	"github.com/yfedoruck/todolist/pkg/pg"
	"net/http"
	"path/filepath"
)

type Router struct {
}

func (r Router) New() {
	db := pg.Postgres{}
	db.Connect()
	defer db.Close()
	db.Tables()

	var loginData = &LoginData{
		"/css/signin.css",
		"Sign in",
		"",
		LoginField{},
	}
	var notesListData = &NotesListData{
		"/css/signin.css",
		"Todo list",
		0,
		nil,
		"",
	}
	var registerData = &RegisterData{
		"/css/signin.css",
		"Sign in",
		RegisterErr{},
		RegisterField{},
	}

	cssDir := env.BasePath() + filepath.FromSlash("/static/css")
	fs := http.FileServer(http.Dir(cssDir))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	http.HandleFunc("/", root)

	http.Handle("/register", &registerHandler{&db, registerData})
	http.Handle("/login", &loginHandler{db, loginData})
	http.Handle("/todolist", todoListHandler(notesListData, db))
	http.Handle("/add", addNoteHandler(notesListData, db))
	http.Handle("/remove", removeTodoHandler(db))
	http.Handle("/logout", logoutHandler(loginData, notesListData, registerData))
}
