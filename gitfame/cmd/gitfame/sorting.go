//go:build !solution

package main

import (
	"sort"
	"strings"
)

func sortByLine(answer []UserData) []UserData {
	sort.SliceStable(answer, func(i, j int) bool {
		if answer[i].lines == answer[j].lines && len(answer[i].commits) != len(answer[j].commits) {
			return len(answer[i].commits) > len(answer[j].commits)
		} else if answer[i].lines == answer[j].lines && len(answer[i].commits) == len(answer[j].commits) && answer[i].files != answer[j].files {
			return answer[i].files > answer[j].files
		} else if answer[i].lines == answer[j].lines && len(answer[i].commits) == len(answer[j].commits) && answer[i].files == answer[j].files {
			return strings.Compare(answer[i].name, answer[j].name) == -1
		} else {
			return answer[i].lines > answer[j].lines
		}
	})
	return answer
}

func sortByCommits(answer []UserData) []UserData {
	sort.SliceStable(answer, func(i, j int) bool {
		if len(answer[i].commits) == len(answer[j].commits) && answer[i].lines != answer[j].lines {
			return answer[i].lines > answer[j].lines
		} else if len(answer[i].commits) == len(answer[j].commits) && answer[i].lines == answer[j].lines && answer[i].files != answer[j].files {
			return answer[i].files > answer[j].files
		} else if len(answer[i].commits) == len(answer[j].commits) && answer[i].lines == answer[j].lines && answer[i].files == answer[j].files {
			return strings.Compare(answer[i].name, answer[j].name) == -1
		} else {
			return len(answer[i].commits) > len(answer[j].commits)
		}
	})
	return answer
}

func sortByFiles(answer []UserData) []UserData {
	sort.SliceStable(answer, func(i, j int) bool {
		if answer[i].files == answer[j].files && answer[i].lines != answer[j].lines {
			return answer[i].lines > answer[j].lines
		} else if answer[i].files == answer[j].files && answer[i].lines == answer[j].lines &&
			len(answer[i].commits) != len(answer[j].commits) {
			return len(answer[i].commits) > len(answer[j].commits)
		} else if answer[i].files == answer[j].files && answer[i].lines == answer[j].lines &&
			len(answer[i].commits) == len(answer[j].commits) {
			return strings.Compare(answer[i].name, answer[j].name) == -1
		} else {
			return answer[i].files > answer[j].files
		}
	})
	return answer
}
