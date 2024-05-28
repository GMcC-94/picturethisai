package main

import (
	"embed"
	"log"
	"log/slog"
	"net/http"
	"os"
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

	router.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(FS))))
	router.Get("/", handler.Make(handler.HandleHomeIndex))
	router.Get("/login", handler.Make(handler.HandleLogin))
	router.Post("/login", handler.Make(handler.HandleLoginPost))

	port := os.Getenv("HTTP_LISTEN_ADDR")
	slog.Info("application running", "port", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initAll() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	return sb.Init()
}
