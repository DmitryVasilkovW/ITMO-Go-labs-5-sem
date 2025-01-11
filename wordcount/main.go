//go:build !solution

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lineCounts := make(map[string]int)

	for _, filename := range os.Args[1:] {
		file, err := os.Open(filename)
		check(err)

		lineCounts = updateMap(*file, lineCounts)
	}

	printCounts(lineCounts)
}

func printCounts(lineCounts map[string]int) {
	for line, count := range lineCounts {
		if count >= 2 {
			fmt.Printf("%d\t%s\n", count, line)
		}
	}
}

func updateMap(file os.File, lineCounts map[string]int) map[string]int {
	scanner := bufio.NewScanner(&file)
	for scanner.Scan() {
		line := scanner.Text()
		lineCounts[line]++
	}

	return lineCounts
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
