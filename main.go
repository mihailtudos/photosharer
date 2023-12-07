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

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.DefaultPostgresConfig()
	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal("Could convert port from .env file")
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = "3af1d7a51d66604a73ea550f8261ebdb"
	cfg.CSRF.Secure = false

	cfg.Server.Address = ":8080"
	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Set up the DB
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Setup services
	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	passwordResetService := &models.PasswordResetService{DB: db}

	emailService := models.NewEmailService(cfg.SMTP)

	// Setup middleware
	umw := controllers.UserMiddleware{SessionService: sessionService}

	csrfMiddleware := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure))

	// Setup controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signin.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "forgot-pw.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "check-your-email.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "reset-pw.gohtml"))

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
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/", usersC.CurrentUser)
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Start server
	fmt.Println("starting server on ", cfg.Server.Address)
	log.Fatal(http.ListenAndServe(cfg.Server.Address, r))
}
