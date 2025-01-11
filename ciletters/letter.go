//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

const firstToolLine = "testtool: copying go.mod, go.sum and .golangci.yml"

//go:embed  text/template/pattern.tmpl
var letterTemplate string

func MakeLetter(n *Notification) (string, error) {
	t, err := makeTemplate()
	if err != nil {
		return "", err
	}

	var message bytes.Buffer
	err = t.Execute(&message, n)
	if err != nil {
		return "", err
	}

	return message.String(), nil
}

func makeTemplate() (*template.Template, error) {
	tmpl, err := template.New("letter").Funcs(template.FuncMap{
		"lastLines": func(s string) []string {
			lines := strings.Split(s, "\n")
			return deleteIndents(lines)
		},
	}).Parse(letterTemplate)
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}

func deleteIndents(lines []string) []string {
	for i, line := range lines {
		if line == firstToolLine {
			lines = lines[i:]
		}
	}

	return lines
}
