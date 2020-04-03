package web

import (
	"github.com/yfedoruck/todolist/pkg/cookie"
	"net/http"
)

func logoutHandler(ld *LoginData, listData *NotesListData, regData *RegisterData) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ld.Error = ""
		listData.Error = ""
		regData.Error = RegisterErr{}
		regData.PreFill = RegisterField{}
		cookie.RemoveCookie(w)
		http.Redirect(w, r, "/login", http.StatusFound)
	})
}
