package web

import (
	"github.com/yfedoruck/todolist/pkg/env"
	"github.com/yfedoruck/todolist/pkg/pg"
	"log"
	"net/http"
)

type Server struct {
	Port string
}

func (s *Server) Start() {
	err := http.ListenAndServe(":"+s.Port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func NewServer(db *pg.Postgres) *Server {
	s := &Server{}
	s.Port = env.Port()
	Router{}.New(db)

	return s
}
