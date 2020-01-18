package validate

import (
	"errors"
	"net/http"
	"strconv"
)

func Note(r *http.Request) error {
	const NoteRule = 50

	note := r.FormValue("note")
	if len(note) > NoteRule {
		return errors.New("max note length = " + strconv.Itoa(NoteRule))
	}
	return nil
}
