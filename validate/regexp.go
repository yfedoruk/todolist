package validate

import (
	"regexp"
)

func Check(str string, rule string) bool {
	var rx = regexp.MustCompile(rule)
	return rx.MatchString(str)
}
