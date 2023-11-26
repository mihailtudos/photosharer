package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello to my awesome site</h1>")
}

func main() {
	http.HandleFunc("/", home)
	fmt.Println("starting server at 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
