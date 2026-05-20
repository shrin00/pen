package main

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shrin00/pen/internal/config"
	"github.com/shrin00/pen/internal/database"
	"github.com/shrin00/pen/internal/repository/postgres"
	"github.com/shrin00/pen/internal/transport/httpx"
)

func main() {
	ctx := context.Background()
	// load the configs
	cfg := config.LoadConfig()
	if cfg.DatabaseURL == "" {
		log.Fatal("database url is required...")
	}

	// connect to database
	db, err := database.NewPostgresPool(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %w", err)
	}
	defer db.Close()
	log.Println("connected to database")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middlePattern)
	router.Use(middleware.Recoverer)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJson(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			httpx.WriteJson(w, http.StatusOK, map[string]string{
				"message": "pong",
			})
		})
	})

	addr := ":" + cfg.Port
	log.Println("server is listening on ", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("failed to run the server: ", err)
	}
}

func middlePattern(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		pattern := chi.RouteContext(r.Context()).RoutePattern()
		log.Println("custome middleware: pattern of the request ", pattern)
	}
	return http.HandlerFunc(fn)
}
