package views

import (
	"bytes"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/mihailtudos/photosharer/context"
	"github.com/mihailtudos/photosharer/models"
	"html/template"
	"io"
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
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("csrfField not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)
	if err != nil {
		return Template{}, fmt.Errorf("parsing fs template: %w", err)
	}

	return Template{htmlTmpl: tpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, r *http.Request, data interface{}, errs ...error) {
	// to avoid race conditions (since template struct is a pointer to a struct)
	// we will copy the templates each time a req comes through
	tpl, err := t.htmlTmpl.Clone()
	if err != nil {
		log.Printf("executing template %v", err.Error())
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// buf is a byte buffer used to write the templates to before writing it to the response
	// it allows to check if there is any error before writing the template and avoiding broken pages
	var buf bytes.Buffer

	// add the csfr token here as we have access to the request
	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return context.User(r.Context())
			},
			"errors": func() []string {
				var errMessages []string
				for _, err := range errs {
					errMessages = append(errMessages, err.Error())
				}
				return errMessages
			},
		},
	)

	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("executing template %v", err.Error())
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(w, &buf)
	if err != nil {
		log.Printf("executing template %v", err.Error())
		http.Error(w, "Failed executing the template", http.StatusInternalServerError)
		return
	}
}
