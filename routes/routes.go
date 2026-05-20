package routes

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shrin00/pen/internal/repository/postgres"
	"github.com/shrin00/pen/internal/transport/handlers"
	"github.com/shrin00/pen/internal/transport/httpx"
	"github.com/shrin00/pen/internal/usecase"
)

func NewRoutes(db *pgxpool.Pool) *chi.Mux {

	userRepo := postgres.NewUserRepository(db)
	authService := usecase.NewAuthService(userRepo)
	AuthHandler := handlers.NewAuthHandler(authService)

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

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", AuthHandler.Register)
			r.Post("/login", AuthHandler.Login)
		})
	})

	return router
}

func middlePattern(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		pattern := chi.RouteContext(r.Context()).RoutePattern()
		log.Println("custome middleware: pattern of the request ", pattern)
	}
	return http.HandlerFunc(fn)
}
