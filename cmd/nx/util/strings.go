package util

import (
	"fmt"
	"strings"
)

func CommaList[T any](things []T, final string) string {
	sb := strings.Builder{}
	for i, thing := range things {
		if i > 0 {
			if i == len(things)-1 {
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
