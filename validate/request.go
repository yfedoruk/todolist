package validate

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
)

func check(str string, rule string) bool {
	var rx = regexp.MustCompile(rule)
	return rx.MatchString(str)
}

func Username(r *http.Request) error {

	const (
		Length = 5
		RegExp = "^[a-zA-Z0-9_]+$"
	)

	var err string

	name := r.FormValue("username")
	if len(name) > Length {
		err = "Max username length = " + strconv.Itoa(Length) + ". "
	}
	if !check(name, RegExp) {
		err = err + "Username must have only digital, alphabetical or underscore symbols. "
	}

	if len(err) > 0 {
		return errors.New(err)
	} else {
		return nil
	}
}

func Email(r *http.Request) error {
	const EmailRule = "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

	email := r.FormValue("email")

	var err string

	if !check(email, EmailRule) {
		err = err + "email is wrong"
	}

	if len(err) > 0 {
		return errors.New(err)
	} else {
		return nil
	}
}

func Note(r *http.Request) error {
	const NoteRule = 50

	note := r.FormValue("note")
	if len(note) > NoteRule {
		return errors.New("max note length = " + strconv.Itoa(NoteRule))
	}
	return nil
}
