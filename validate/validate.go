package validate

import (
	"net/http"
	"regexp"
)

func Validate(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f(w, r)
	}
}

func Check(str string, rule string) bool {
	var rx = regexp.MustCompile(rule)
	return rx.MatchString(str)
}
