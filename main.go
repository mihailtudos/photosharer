package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"html/template"
	"log"
	"net/http"
	"path"
)

type Router struct{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := path.Join("templates", "home.gohtml")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("parsing template %s", tmplPath)
		http.Error(w, "Failed parsing the template", http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Printf("executing template %s", tmplPath)
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	fmt.Fprintln(w, "<h1>Welcome to the contact page</h1>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	fmt.Fprintln(w, "<h1>Welcome to the FAQ page</h1>")
}

func main() {
	r := chi.NewRouter()
	r.Get("/", homeHandler)

	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	fmt.Println("starting server at 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
