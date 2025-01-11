//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	urls := os.Args[1:]

	for _, url := range urls {
		resp, err := http.Get(url)
		stopAndShowError(err, url)

		body, err := io.ReadAll(resp.Body)
		stopAndShowError(err, url)

		fmt.Printf("%s", body)
	}
}

func stopAndShowError(err error, url string) {
	code := handleErrorAndGetExitCode(err, url)
	if code != -1 {
		os.Exit(code)
	}
}

func handleErrorAndGetExitCode(ex error, url string) int {
	if ex != nil {
		_, err := fmt.Fprintf(os.Stderr, "fetch: %s: %v\n", url, ex)
		if err != nil {
			return 0
		}
		return 1
	}
	return -1
}
