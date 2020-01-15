package middleware

import (
	"errors"
	"net/http"
)

const (
	EMAIL_RULE = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

func Email(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	if Check(email, EMAIL_RULE) != false {
		return errors.New("email is wrong")
	}
	return nil
}
