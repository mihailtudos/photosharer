package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Template struct {
	htmlTmpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(filePath string) (Template, error) {
	t, err := template.ParseFiles(filePath)
	if err != nil {
		return Template{}, fmt.Errorf("parsing template: %w", err)
	}

	return Template{htmlTmpl: t}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTmpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template %v", err.Error())
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}
}
