package env

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

	//baseDir, _ := os.Getwd()
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Caller error")
	}
	envDir := filepath.Dir(b)
	pkgDir := filepath.Dir(envDir)
	appDir := filepath.Dir(pkgDir)

	return appDir
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
