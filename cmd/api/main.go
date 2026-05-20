package main

import (
	"log"
	"net/http"

	"github.com/shrin00/pen/internal/config"
	"github.com/shrin00/pen/internal/database"
	"github.com/shrin00/pen/routes"
)

func main() {
	// ctx := context.Background()
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

	router := routes.NewRoutes(db)

	addr := ":" + cfg.Port
	log.Println("server is listening on ", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal("failed to run the server: ", err)
	}
}
