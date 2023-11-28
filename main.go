package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type Router struct{}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello to my awesome site</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	fmt.Fprintln(w, "<h1>Welcome to the contact page</h1>")
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method, r.URL.Path)
	fmt.Fprintln(w, "<h1>Welcome to the FAQ page</h1>")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	user := r.Context().Value("user")

	fmt.Fprintln(w, fmt.Sprintf("<h1>Showing artile with id %s and user is %s</h1>", id, user))
}

func MyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "user", "123")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	r := chi.NewRouter()

	r.Use(MyMiddleware)
	r.Get("/", homeHandler)

	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	r.Get("/article/{id:[0-9]+}", testHandler)

	fmt.Println("starting server at 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
