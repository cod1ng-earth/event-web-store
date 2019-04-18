package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Hello, world!")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8000")
		fmt.Fprintf(w, "Hello, you have requested: %s %d", r.URL.Path, t.UnixNano())
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
