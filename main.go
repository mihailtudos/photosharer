package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mihailtudos/photosharer/views"
	"log"
	"net/http"
	"path"
)

type Router struct{}

func render(w http.ResponseWriter, filePath string, data interface{}) {
	t, err := views.Parse(filePath)
	if err != nil {
		log.Printf("executing template %v", err.Error())
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}

	t.Execute(w, data)
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := path.Join("templates", "home.gohtml")
	render(w, tmplPath, nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := path.Join("templates", "contact.gohtml")
	render(w, tmplPath, nil)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	tmplPath := path.Join("templates", "faq.gohtml")
	render(w, tmplPath, nil)
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
