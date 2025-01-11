//go:build !solution

package main

import "path"

func excludes(request RepoFlags, fileName string) bool {
	if len(request.exclude) == 0 {
		return false
	}

	for _, pattern := range request.exclude {
		match, err := path.Match(pattern, fileName)
		if err != nil {
			return true
		}
		if match {
			return true
		}
	}

	return false
}

func restriced(request RepoFlags, fileName string) bool {
	if len(request.restrictTo) == 0 {
		return false
	}

	for _, pattern := range request.restrictTo {
		match, err := path.Match(pattern, fileName)
		if err != nil {
			return true
		}
		if match {
			return false
		}
	}

	return true
}
