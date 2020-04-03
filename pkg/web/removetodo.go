package web

import (
	"github.com/yfedoruck/todolist/pkg/pg"
	"github.com/yfedoruck/todolist/pkg/resp"
	"log"
	"net/http"
	"strconv"
)

func removeTodoHandler(db pg.Postgres) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			err := r.ParseForm()
			resp.Check(err)

			if len(r.Form["id"]) == 0 {
				panic("id not exists")
			}

			if _, err := r.Cookie("auth"); err != nil {
				log.Println("user id not found")
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}

			ok, err := strconv.Atoi(r.PostFormValue("id"))
			resp.Check(err)

			db.RemoveNote(ok)
			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	})
}
