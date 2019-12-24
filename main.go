package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
)
var ErrNotImplemented = errors.New("not implemented")

func main() {
	addr := ":8080"
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {})
	r.Post("/token", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {})
	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {})

	log.Fatalln(http.ListenAndServe(addr, r))
}

func connect() (*sqlx.DB, error) {
	return nil, ErrNotImplemented
}

func storeToken(username, authToken, apiKey string) error {
	return nil
}
func loadToken(username string) (string, string, error) {
	return "", "", nil
}
