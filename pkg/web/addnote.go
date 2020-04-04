package web

import (
	"github.com/yfedoruck/todolist/pkg/cookie"
	"github.com/yfedoruck/todolist/pkg/pg"
	"github.com/yfedoruck/todolist/pkg/resp"
	"github.com/yfedoruck/todolist/pkg/validate"
	"net/http"
)

func addNoteHandler(notes *NotesListData, db *pg.Postgres) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth, err := r.Cookie("auth")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			resp.Check(err)

			if len(r.Form["note"]) == 0 {
				panic("note not exists")
			}

			err = validate.Note(r)
			if err != nil {
				notes.Error = err.Error()
				http.Redirect(w, r, "/todolist", http.StatusFound)
				return
			} else {
				notes.Error = ""
			}

			user := cookie.Cookie{}
			user.Decode(auth.Value)
			db.AddNote(user.Id, r.PostFormValue("note"))

			http.Redirect(w, r, "/todolist", http.StatusFound)
		}
	})
}
