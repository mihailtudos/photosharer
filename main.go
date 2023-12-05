package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
	"github.com/mihailtudos/photosharer/controllers"
	"github.com/mihailtudos/photosharer/migrations"
	"github.com/mihailtudos/photosharer/models"
	"github.com/mihailtudos/photosharer/templates"
	"github.com/mihailtudos/photosharer/views"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal("Could convert port from .env file")
	}
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	// Set up the DB
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	models.NewEmailService(models.SMTPConfig{Host: host, Port: port, Username: username, Password: password})

	// Setup middleware
	umw := controllers.UserMiddleware{SessionService: &sessionService}

	csrfKey := "3af1d7a51d66604a73ea550f8261ebdb"
	csrfMiddleware := csrf.Protect([]byte(csrfKey), csrf.Secure(false))

	// Setup controllers
	usersC := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signin.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "forgot-pw.gohtml"))

	//setup router and routes
	r := chi.NewRouter()
	r.Use(csrfMiddleware)
	r.Use(umw.SetUser)

	tpl := views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "home.gohtml"))
	r.Get("/", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "contact.gohtml"))
	r.Get("/contact", controllers.StaticHandler(tpl))

	tpl = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "faq.gohtml"))
	r.Get("/faq", controllers.FAQ(tpl))

	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/users", usersC.Create)
	r.Post("/signin", usersC.Authenticate)
	r.Post("/signout", usersC.SignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Start server
	fmt.Println("starting server at 8080...")
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
