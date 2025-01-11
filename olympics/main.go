package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var athletes []Information

func main() {
	port := flag.String("port", "80", "http server port")
	jsonP := flag.String("data", "./olympics/testdata/olympicWinners.json", "json path")
	flag.Parse()

	jF, err := os.Open(*jsonP)
	if err != nil {
		log.Fatal("file read err")
	}

	byteVal, _ := io.ReadAll(jF)
	err = jF.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(byteVal, &athletes)

	if err != nil {
		log.Fatal("json parse err")
	}

	http.HandleFunc("/athlete-info", Athletes)
	http.HandleFunc("/top-athletes-in-sport", TopAthletes)
	http.HandleFunc("/top-countries-in-year", TopCountries)

	host := fmt.Sprintf(":%s", *port)
	log.Fatal(http.ListenAndServe(host, nil))
}
