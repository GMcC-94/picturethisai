package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
	"picturethisai/db"
	"picturethisai/handler"
	"picturethisai/pkg/sb"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

//go:embed public
var FS embed.FS

func main() {
	if err := initAll(); err != nil {
		log.Fatal(err)
	}

	router := chi.NewMux()
	router.Use(handler.WithUser)

	// Public Routes
	router.Handle("/*", public())
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleLoginIndex))
	router.Get("/login/provider/google", handler.Make(handler.HandleLoginWithGoogle))
	router.Get("/signup", handler.Make(handler.HandleSignupIndex))
	router.Get("/auth/callback", handler.Make(handler.HandleAuthCallback))
	router.Get("/account/setup", handler.Make(handler.HandleAccountSetupIndex))

	router.Post("/logout", handler.Make(handler.HandleLogoutCreate))
	router.Post("/login", handler.Make(handler.HandleLoginCreate))
	router.Post("/signup", handler.Make(handler.HandleSignupCreate))
	router.Post("/account/setup", handler.Make(handler.HandleAccountSetupCreate))

	// Auth routes
	router.Group(func(auth chi.Router) {
		auth.Use(handler.WithAccountSetup)
		auth.Get("/settings", handler.Make(handler.HandleSettingsIndex))
		auth.Put("/settings/account/profile", handler.Make(handler.HandleSettingsUsernameUpdate))
	})

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initAll() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	if err := db.Init(); err != nil {
		return err
	}

	return sb.Init()
}
