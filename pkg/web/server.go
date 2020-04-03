package web

import (
	"github.com/yfedoruck/todolist/pkg/env"
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

func NewServer() *Server {
	s := &Server{}
	s.Port = env.Port()
	Router{}.New()

	return s
}
