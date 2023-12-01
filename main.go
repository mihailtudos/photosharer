package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mihailtudos/photosharer/controllers"
	"github.com/mihailtudos/photosharer/models"
	"github.com/mihailtudos/photosharer/templates"
	"github.com/mihailtudos/photosharer/views"
	"log"
	"net/http"
)

func main() {
	r := chi.NewRouter()
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userService := models.UserService{
		DB: db,
	}

	tpl := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	usersC := controllers.Users{
		UserService: &userService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signin.gohtml"))
	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/users", usersC.Create)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	fmt.Println("starting server at 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
