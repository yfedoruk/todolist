package web

import (
	"github.com/yfedoruck/todolist/pkg/cookie"
	"github.com/yfedoruck/todolist/pkg/pg"
	"net/http"
)

func todoListHandler(data *NotesListData, db pg.Postgres) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, err := r.Cookie("auth")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		user := cookie.Cookie{}
		user.Decode(auth.Value)
		data.UserId = user.Id
		data.TodoList = db.TodoListData(user.Id)
		renderTemplate(w, "todolist", data)
	})
}
