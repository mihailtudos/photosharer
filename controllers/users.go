package controllers

import (
	"fmt"
	"net/http"
)

// Users struct provides functionality needed for the users controller
// Decoupling the views from the controllers by creating a Template interface makes the controllers
// agnostic to how date is returned to the user. Thus allowing the same controller to power an API and/or HTML
// another benefit of decoupling here is breaking circular dependencies
type Users struct {
	Templates struct {
		New Template
	}
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "email %s", r.FormValue("email"))
	fmt.Fprintf(w, "name %s", r.FormValue("name"))
}
