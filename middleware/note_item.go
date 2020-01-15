package middleware

import (
	"errors"
	"net/http"
	"strconv"
)

const (
	NOTE_RULE = 50
)

func Note(w http.ResponseWriter, r *http.Request) error {
	note := r.FormValue("note")
	if len(note) > NOTE_RULE {
		return errors.New("max note length = " + strconv.Itoa(NOTE_RULE))
	}
	return nil
}
