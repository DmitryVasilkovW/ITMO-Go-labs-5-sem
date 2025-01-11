//go:build !solution

package speller

import "strings"

var (
	ones = []string{
		"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine",
		"ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen",
	}

	tens = []string{
		"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety",
	}

	powers = []struct {
		value int64
		name  string
	}{
		{1_000_000_000_000, "trillion"},
		{1_000_000_000, "billion"},
		{1_000_000, "million"},
		{1_000, "thousand"},
		{100, "hundred"},
	}
)

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}

	if n < 0 {
		return "minus " + Spell(-n)
	}

	resultString := addAllNumbers(n)
	return strings.Join(resultString, " ")
}

func addAllNumbers(n int64) []string {
	var parts []string

	for _, power := range powers {
		if n >= power.value {
			parts = append(parts, spellThreeDigits(n/power.value)+" "+power.name)
			n %= power.value
		}
	}

	if n > 0 {
		parts = append(parts, spellThreeDigits(n))
	}

	return parts
}

func spellThreeDigits(n int64) string {
	if n == 0 {
		return ""
	}

	var parts []string
	parts = *tryToAddHundreds(&n, &parts)
	parts = tryToAddTensAndOnes(n, parts)

	return strings.Join(parts, " ")
}

func tryToAddHundreds(n *int64, parts *[]string) *[]string {
	if *n >= 100 {
		*parts = append(*parts, ones[*n/100]+" hundred")
		*n %= 100
	}

	return parts
}

func tryToAddTensAndOnes(n int64, parts []string) []string {
	if n >= 20 {
		parts = append(parts, tens[n/10])
		if n%10 != 0 {
			parts[len(parts)-1] += "-" + ones[n%10]
		}
	} else if n > 0 {
		parts = append(parts, ones[n])
	}

	return parts
}
