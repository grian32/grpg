package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	log.Printf("serving ./assets on port 4022")
	log.Fatal(http.ListenAndServe(":4022", nil))
}
