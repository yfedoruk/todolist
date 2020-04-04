package main

import (
	"github.com/yfedoruck/todolist/pkg/pg"
	"github.com/yfedoruck/todolist/pkg/web"
)

type App struct {
	server *web.Server
	db     *pg.Postgres
}

func (a *App) Init() {
	db := &pg.Postgres{}
	db.Connect()
	a.db = db

	a.server = web.NewServer(a.db)
}

func (a *App) Run() {
	a.db.Connect()
	defer a.db.Close()
	a.db.Tables()
	a.server.Start()
}
