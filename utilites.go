package main

import (
	"log"
	"path/filepath"
	"runtime"
)

var baseDir string

func basePath() string {
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
