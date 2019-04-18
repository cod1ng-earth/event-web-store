package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Hello, world.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, you have requested: %s\n", r.URL.Path)
	})

	http.ListenAndServe(":8080", nil)
}
