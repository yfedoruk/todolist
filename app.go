package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var baseDir string

func BasePath() string {
	if baseDir != "" {
		return baseDir
	}
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Caller error")
	}

	baseDir = filepath.Dir(b)
	return baseDir
}

func Env() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
}

