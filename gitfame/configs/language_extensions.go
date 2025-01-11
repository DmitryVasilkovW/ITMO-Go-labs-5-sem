//go:build !solution

package configs

import (
	"encoding/json"
	"os"
	"strings"
)

type Language struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Extensions []string `json:"extensions"`
}

var mapOfLanguages = make(map[string][]string)

func GetExts(lang string) []string {
	if len(mapOfLanguages) == 0 {
		var languages []Language
		content, _ := os.ReadFile("../../configs/language_extensions.json")
		_ = json.Unmarshal(content, &languages)

		// fmt.Fprintln(os.Stderr, content)

		for _, lang := range languages {
			mapOfLanguages[strings.ToLower(lang.Name)] = lang.Extensions
			//fmt.Fprintln(os.Stderr, lang.name)
		}
	}
	return mapOfLanguages[strings.ToLower(lang)]
}
