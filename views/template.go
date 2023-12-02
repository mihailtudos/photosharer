package views

import (
	"fmt"
	"github.com/gorilla/csrf"
	"html/template"
	"io/fs"
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

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// Template functions need to be added before the templates are parsed
	tpl := template.New(patterns[0])
	tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return `<!-- TODO: implement the csrfField -->`
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing fs template: %w", err)
	}

	return Template{htmlTmpl: tpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}) {
	// to avoid race conditions (since template struct is a pointer to a struct)
	// we will copy the templates each time a req comes through
	tpl, err := t.htmlTmpl.Clone()
	if err != nil {
		log.Printf("executing template %v", err.Error())
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// add the csfr token here as we have access to the request
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
		},
	)

	err = tpl.Execute(w, data)
	if err != nil {
		log.Printf("executing template %v", err.Error())
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}
}
