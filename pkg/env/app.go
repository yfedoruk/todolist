package env

import (
	"os"
)

var baseDir string

func BasePath() string {
	if baseDir != "" {
		return baseDir
	}

	baseDir, _ := os.Getwd()
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


