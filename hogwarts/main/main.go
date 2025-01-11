package main

import (
	"fmt"
	"gitlab.com/slon/shad-go/hogwarts"
)

func main() {
	var linearScience = map[string][]string{
		"1": {"0"},
		"2": {"1"},
		"3": {"2"},
		"4": {"3"},
		"5": {"4"},
		"6": {"5"},
		"7": {"6"},
		"8": {"7"},
		"9": {"8"},
	}

	fmt.Println(hogwarts.GetCourseList(linearScience))
}
