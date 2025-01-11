//go:build !solution

package main

import (
	"gitlab.com/slon/shad-go/gitfame/configs"
	"path/filepath"
)

func extCheck(request RepoFlags, fileName string) bool {
	if len(request.extensions) == 0 {
		return false
	}

	fileExtension := filepath.Ext(fileName)

	for _, ext := range request.extensions {
		if ext == fileExtension {
			return false
		}
	}

	return true
}

func langCheck(request RepoFlags, fileName string) bool {
	if len(request.languages) == 0 {
		return false
	}

	fileExtension := filepath.Ext(fileName)

	for _, lang := range request.languages {
		langExts := configs.GetExts(lang)

		for _, ext := range langExts {
			if ext == fileExtension {
				return false
			}
		}
	}

	return true
}
