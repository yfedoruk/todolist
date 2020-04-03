package web

import "github.com/yfedoruck/todolist/pkg/pg"

type NotesListData struct {
	Css      string
	Title    string
	UserId   int
	TodoList []pg.Todo
	Error    string
}
