package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var (
	KeyToURL        map[string]string
	URLToKey        map[string]string
	randomGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type URL struct {
	Value string `json:"url"`
}

type Response struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

func redirectToURL(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/go/"):]

	url, exists := KeyToURL[key]
	if !exists {
		http.Error(w, "invalid key", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func createShortURL(w http.ResponseWriter, r *http.Request) {
	var url URL
	if err := json.NewDecoder(r.Body).Decode(&url); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if key, exists := URLToKey[url.Value]; exists {
		writeJSONResponse(w, http.StatusOK, Response{URL: url.Value, Key: key})
		return
	}

	key := generateKey()
	KeyToURL[key] = url.Value
	URLToKey[url.Value] = key

	writeJSONResponse(w, http.StatusOK, Response{URL: url.Value, Key: key})
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func generateKey() string {
	return strconv.Itoa(randomGenerator.Intn(1e9))
}

func main() {
	port := flag.Int("port", 8080, "port string")
	flag.Parse()

	initializeMaps()
	setupRoutes()

	address := "localhost:" + strconv.Itoa(*port)
	log.Printf("Starting server on %s", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func initializeMaps() {
	KeyToURL = make(map[string]string)
	URLToKey = make(map[string]string)
}

func setupRoutes() {
	http.HandleFunc("/shorten", createShortURL)
	http.HandleFunc("/go/", redirectToURL)
}
