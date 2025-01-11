//go:build !solution

package reverse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func Reverse(input string) string {
	var sb strings.Builder
	sb.Grow(len(input))

	for i := len(input); i > 0; {
		r, size := utf8.DecodeLastRuneInString(input[:i])
		addRune(r, size, &sb)

		i -= size
	}

	return sb.String()
}

func addRune(r rune, size int, sb *strings.Builder) {
	if r == utf8.RuneError && size == 1 {
		sb.WriteRune(unicode.ReplacementChar)
		return
	}

	sb.WriteRune(r)
}
