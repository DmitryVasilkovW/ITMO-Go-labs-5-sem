//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	timeOfStart := time.Now()
	ch := make(chan string)
	countOfUrls := len(os.Args[1:])

	runAllFetchers(ch)

	showAllData(ch, countOfUrls)
	showElapsedTime(timeOfStart)
}

func runAllFetchers(ch chan string) {
	for _, url := range os.Args[1:] {
		go fetch(url, ch)
	}
}

func showAllData(ch chan string, countOfUrls int) {
	for range countOfUrls {
		fmt.Println(<-ch)
	}
}

func showElapsedTime(start time.Time) {
	fmt.Println(getMessageForElapsedTime(start))
}

func getMessageForElapsedTime(start time.Time) string {
	secs := time.Since(start).Seconds()
	return fmt.Sprintf("%.2fs elapsed\n", secs)
}

func fetch(url string, ch chan<- string) {
	timeOfStart := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		ch <- getErrorMessageForGetRequest(err, url)
		return
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		ch <- getErrorMessageForGetRequest(err, url)
		return
	}
	ex := resp.Body.Close()
	if ex != nil {
		return
	}

	ch <- getMessageForGetRequest(url, data, timeOfStart)
}

func getErrorMessageForGetRequest(err error, url string) string {
	return fmt.Sprintf("Get %s: %v", url, err)
}

func getMessageForGetRequest(url string, data []byte, timeOfStart time.Time) string {
	secs := time.Since(timeOfStart).Seconds()
	return fmt.Sprintf("%.2fs  %7d  %s", secs, data, url)
}
