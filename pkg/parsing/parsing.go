package parsing

import (
	"strings"
)

func CleanupTitle(s string) string {
	s = strings.TrimPrefix(s, "[[")
	s = strings.TrimSuffix(s, "]]")

	s = strings.Split(s, "|")[0]
	s = strings.Split(s, "#")[0]

	t := strings.Split(s, ":")
	if len(t) != 1 {
		s = t[1]
	}

	s = strings.TrimLeft(s, " ")
	s = strings.TrimRight(s, " ")

	s = strings.ToLower(s)

	return s
}
