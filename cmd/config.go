package main

import "os"

func port() string {
	p := os.Getenv("PORT")
	if p == "" {
		p = "5000"
	}
	return p
}
