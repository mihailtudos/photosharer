package controllers

import (
	"html/template"
	"net/http"
)

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Where are you located?",
			Answer:   "We are based in London, UK",
		},
		{
			Question: "How can we contact you?",
			Answer:   `You can reach to us via phone <a href="tel:+44 7982193932">+44 7982193932</a> of email <a href="mailto:contact@renect.com">contact@renect.com</a> `,
		},
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, r, questions)
	}
}
