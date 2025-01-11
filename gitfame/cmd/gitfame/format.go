//go:build !solution

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"text/tabwriter"
)

func writeTabular(answer []UserData) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintln(w, "Name\tLines\tCommits\tFiles")
	for _, k := range answer {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d", k.name, k.lines, len(k.commits), k.files)
		fmt.Fprintln(w)
	}
	w.Flush()
	return nil
}

func writeCSV(answer []UserData) error {
	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()
	err := csvWriter.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		return err
	}
	for _, user := range answer {
		err := csvWriter.Write([]string{user.name, fmt.Sprint(user.lines), fmt.Sprint(len(user.commits)), fmt.Sprint(user.files)})
		if err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return nil
}

func writeJSON(answer []UserData) error {
	fmt.Printf("%s", "[")
	for i, k := range answer {
		fmt.Printf("{\"name\":\"%s\",\"lines\":%d,\"commits\":%d,\"files\":%d}", k.name, k.lines, len(k.commits), k.files)
		if i != len(answer)-1 {
			fmt.Printf(",")
		}
	}
	fmt.Println("]")
	return nil
}

func writeJSONLines(answer []UserData) error {
	for _, k := range answer {
		fmt.Printf("{\"name\":\"%s\",\"lines\":%d,\"commits\":%d,\"files\":%d}\n", k.name, k.lines, len(k.commits), k.files)
	}
	return nil
}
