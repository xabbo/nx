package util

import (
	"fmt"
	"strings"
)

func CommaList[T any](things []T, final string) string {
	sb := strings.Builder{}
	for i, thing := range things {
		if i > 0 {
			if i == len(things)-1 && final != "" {
				sb.WriteRune(' ')
				sb.WriteString(final)
				sb.WriteRune(' ')
			} else {
				sb.WriteString(", ")
			}
		}
		sb.WriteString(fmt.Sprint(thing))
	}
	return sb.String()
}

func Pluralize(n int, word string, plural string) string {
	if n != 1 {
		word += plural
	}
	return fmt.Sprintf("%d %s", n, word)
}
