package env

import (
	"github.com/yfedoruck/todolist/pkg/resp"
	"os"
)

var baseDir string

func BasePath() string {
	if baseDir != "" {
		return baseDir
	}

	baseDir, err := os.Getwd()
	resp.Check(err)

	return baseDir
}

func Domain() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
}

func Port() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "5000"
	}
	return p
}
