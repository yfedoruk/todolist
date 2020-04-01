package main

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

func Env() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
}
