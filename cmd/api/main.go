package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/shrin00/pen/internal/config"
	"github.com/shrin00/pen/internal/transport/httpx"
)

func main() {

	cfg := config.LoadConfig()

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
