package controllers

import (
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/mihailtudos/photosharer/models"
	"html/template"
	"log"
	"net/http"
)

// Users struct provides functionality needed for the users controller
// Decoupling the views from the controllers by creating a Template interface makes the controllers
// agnostic to how date is returned to the user. Thus allowing the same controller to power an API and/or HTML
// another benefit of decoupling here is breaking circular dependencies
type Users struct {
	Templates struct {
		New    Template
		SignIn Template
	}
	UserService *models.UserService
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email string
	}

	data.Email = r.FormValue("email")
	u.Templates.New.Execute(w, data)
}

func (u Users) SignIn(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email     string
		CSRFField template.HTML
	}

	data.Email = r.FormValue("email")
	data.CSRFField = csrf.TemplateField(r)
	u.Templates.SignIn.Execute(w, data)
}

func (u Users) Authenticate(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Email    string
		Password string
	}

	data.Email = r.FormValue("email")
	data.Password = r.FormValue("password")

	user, err := u.UserService.Authenticate(data.Email, data.Password)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "email", Value: user.Email, Path: "/", HttpOnly: true}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "Authenticated %+v", user)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := u.UserService.Create(email, password)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%#v", user)
}

func (u Users) CurrentUser(w http.ResponseWriter, r *http.Request) {
	email, err := r.Cookie("email")
	if err != nil {
		fmt.Fprintf(w, "This eamil cookie could not be read")
	}

	fmt.Fprintf(w, "email cookie: %s\n", email.Value)
}
