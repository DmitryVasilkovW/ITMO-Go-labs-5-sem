//go:build !solution

package fileleak

import (
	"io/fs"
	"os"
	"reflect"
)

const dir = "/proc/self/fd"

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func getInfoAboutFiles() []fs.FileInfo {
	entries, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	return getInfo(entries)
}

func getInfo(entries []os.DirEntry) []fs.FileInfo {
	infoList := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()

		if err != nil {
			continue
		}

		infoList = append(infoList, info)
	}

	return infoList
}

func VerifyNone(t testingT) {
	before := getInfoAboutFiles()
	t.Cleanup(func() {
		getAfter(t, before)
	})
}

func getAfter(t testingT, before []os.FileInfo) {
	after := getInfoAboutFiles()
	if !reflect.DeepEqual(after, before) {
		t.Errorf("Leaks detected. Before: %v, After: %v", before, after)
	}
}
