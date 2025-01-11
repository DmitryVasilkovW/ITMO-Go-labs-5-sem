//go:build !solution

package varfmt

import (
	"fmt"
	"strconv"
	"strings"
)

func Sprintf(format string, args ...interface{}) string {
	var sb strings.Builder
	allocateMemoryForAllArgs(&sb, &args, &format)

	iteratorForArgs := 0
	for i := 0; i < len(format); i++ {
		if format[i] == '{' {
			parserForBraces(&sb, format, &i, &iteratorForArgs, &args)
		} else {
			sb.WriteByte(format[i])
		}
	}

	return sb.String()
}

func allocateMemoryForAllArgs(sb *strings.Builder, args *[]interface{}, format *string) {
	totalLength := len(*format)
	for _, arg := range *args {
		switch v := arg.(type) {
		case int:
			totalLength += len(strconv.Itoa(v))
		case float64:
			totalLength += len(strconv.FormatFloat(v, 'f', -1, 64))
		case string:
			totalLength += len(v)
		default:
			totalLength += len(fmt.Sprint(v))
		}
	}
	sb.Grow(totalLength)
}

func parserForBraces(sb *strings.Builder, format string, i, iteratorForArgs *int, args *[]interface{}) {
	indexForRightBrace := *i + 1
	for indexForRightBrace < len(format) && format[indexForRightBrace] != '}' {
		indexForRightBrace++
	}
	if indexForRightBrace == len(format) {
		panic("incorrect format")
	}

	addArg(sb, format, i, *iteratorForArgs, indexForRightBrace, args)
	*i = indexForRightBrace
	*iteratorForArgs++
}

func addArg(sb *strings.Builder, format string, i *int, iteratorForArgs, indexForRightBrace int, args *[]interface{}) {
	insideBraces := format[*i+1 : indexForRightBrace]
	if insideBraces == "" {
		addArgByIterator(iteratorForArgs, sb, *args)
		return
	}

	addArgByIndex(&insideBraces, sb, args)
}

func addArgByIndex(insideBraces *string, sb *strings.Builder, args *[]interface{}) {
	argIndex, err := strconv.Atoi(*insideBraces)
	if err != nil || argIndex < 0 || argIndex >= len(*args) {
		panic("incorrect format")
	}

	addArgByIterator(argIndex, sb, *args)
}

func addArgByIterator(index int, sb *strings.Builder, args []interface{}) {
	switch v := args[index].(type) {
	case int:
		sb.WriteString(strconv.Itoa(v))
	case float64:
		sb.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	case string:
		sb.WriteString(v)
	default:
		sb.WriteString(fmt.Sprint(v))
	}
}
