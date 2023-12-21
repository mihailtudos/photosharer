package main

import (
	"fmt"
	"github.com/mihailtudos/photosharer/controllers"
	"golang.org/x/oauth2"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
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
	OAuthProviders map[string]*oauth2.Config
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}

	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		Password: os.Getenv("PSQL_PASSWORD"),
		Database: os.Getenv("PSQL_DATABASE"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL condig provided")
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	cfg.SMTP.Port, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	cfg.OAuthProviders = make(map[string]*oauth2.Config)
	cfg.OAuthProviders["dropbox"] = &oauth2.Config{
		ClientID:     os.Getenv("DROPBOX_APP_ID"),
		ClientSecret: os.Getenv("DROPBOX_APP_SECRET"),
		Scopes:       []string{"files.content.read", "files.metadata.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}

	return cfg, nil
}

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	err = run(cfg)
	if err != nil {
		panic(err)
	}
}

func run(cfg config) error {
	// Set up the DB
	db, err := models.Open(cfg.PSQL)
	if err != nil {
		return err
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		return err
	}

	// Setup services
	userService := &models.UserService{DB: db}
	sessionService := &models.SessionService{DB: db}
	passwordResetService := &models.PasswordResetService{DB: db}
	galleryService := &models.GalleryService{DB: db}

	emailService := models.NewEmailService(cfg.SMTP)

	// Setup middleware
	umw := controllers.UserMiddleware{SessionService: sessionService}

	csrfMiddleware := csrf.Protect([]byte(cfg.CSRF.Key), csrf.Secure(cfg.CSRF.Secure), csrf.Path("/"))

	// Setup controllers
	usersC := controllers.Users{
		UserService:          userService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
		EmailService:         emailService,
	}

	//auth controllers
	oauthController := controllers.OAuth{
		ProviderConfigs: cfg.OAuthProviders,
	}

	usersC.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signup.gohtml"))
	usersC.Templates.SignIn = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "signin.gohtml"))
	usersC.Templates.ForgotPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "forgot-pw.gohtml"))
	usersC.Templates.CheckYourEmail = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "check-your-email.gohtml"))
	usersC.Templates.ResetPassword = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "reset-pw.gohtml"))

	galleriesController := controllers.Galleries{
		GalleryService: galleryService,
	}

	galleriesController.Templates.New = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "galleries/new.gohtml"))
	galleriesController.Templates.Edit = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "galleries/edit.gohtml"))
	galleriesController.Templates.Index = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "galleries/index.gohtml"))
	galleriesController.Templates.Show = views.Must(views.ParseFS(templates.FS, "layout.gohtml", "shared/navbar.gohtml", "shared/footer.gohtml", "galleries/show.gohtml"))

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

	// Users related routes
	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/users", usersC.Create)
	r.Post("/signin", usersC.Authenticate)
	r.Post("/signout", usersC.SignOut)
	r.Get("/forgot-pw", usersC.ForgotPassword)
	r.Post("/forgot-pw", usersC.ProcessForgotPassword)
	r.Get("/reset-pw", usersC.ResetPassword)
	r.Post("/reset-pw", usersC.ProcessResetPassword)

	// Galleries related routes
	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleriesController.Show)
		r.Get("/{id}/images/{filename}", galleriesController.Image)

		r.Group(func(r chi.Router) {
			r.Use(umw.RequireUser)
			r.Get("/", galleriesController.Index)
			r.Get("/new", galleriesController.New)
			r.Post("/", galleriesController.Create)
			r.Post("/{id}", galleriesController.Update)
			r.Post("/{id}/images", galleriesController.UploadImages)
			r.Get("/{id}/edit", galleriesController.Edit)
			r.Post("/{id}/delete", galleriesController.Delete)
			r.Post("/{id}/images/{filename}/delete", galleriesController.DeleteImage)
		})
	})

	r.Route("/oauth/{provider}", func(r chi.Router) {
		r.Use(umw.RequireUser)
		r.Get("/connect", oauthController.Connect)
		r.Get("/callback", oauthController.Callback)
	})

	assetsHandler := http.FileServer(http.Dir("assets"))
	r.Get("/assets/*", http.StripPrefix("/assets", assetsHandler).ServeHTTP)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	// Start server
	fmt.Println("starting server on ", cfg.Server.Address)
	return http.ListenAndServe(cfg.Server.Address, r)
}
