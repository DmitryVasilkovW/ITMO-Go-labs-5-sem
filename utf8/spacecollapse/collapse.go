//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var builder strings.Builder
	builder.Grow(len(input))
	isFirstSpace := true

	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		addRune(r, size, &isFirstSpace, &builder)

		i += size
	}

	return builder.String()
}

func addRune(r rune, size int, isFirstSpace *bool, sb *strings.Builder) {
	if r == utf8.RuneError && size == 1 {
		writeRuneAndSetFlag(sb, unicode.ReplacementChar, isFirstSpace, true)
		return
	}

	if unicode.IsSpace(r) && *isFirstSpace {
		writeRuneAndSetFlag(sb, ' ', isFirstSpace, false)
	} else if !unicode.IsSpace(r) {
		writeRuneAndSetFlag(sb, r, isFirstSpace, true)
	}
}

func writeRuneAndSetFlag(sb *strings.Builder, r rune, isFirstSpace *bool, flag bool) {
	sb.WriteRune(r)
	*isFirstSpace = flag
}
