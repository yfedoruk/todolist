package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type cookie struct {
	Name string
	Id   int
}

func (c cookie) encode() string {
	js, err := json.Marshal(c)
	check(err)

	return base64.StdEncoding.EncodeToString(js)
}

func (c *cookie) decode(arg string) {
	js, err := base64.StdEncoding.DecodeString(arg)
	check(err)

	err = json.Unmarshal(js, &c)
	check(err)
}

func (c cookie) set(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:  "auth",
		Value: c.encode(),
		Path:  "/",
	})
}

func removeCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "auth",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
