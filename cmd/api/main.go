package main

import (
	"log"
	"net/http"
	"os"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"api-avis-clients/internal/config"
	"api-avis-clients/internal/handler"
	"api-avis-clients/internal/repository"
)

func main() {
	// Initialisation de la base de données
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données: %v", err)
	}
	defer config.CloseDB()

	// Initialisation du repository
	avisRepo := repository.NewAvisRepository(db)

	// Initialisation du handler
	avisHandler := handler.NewAvisHandler(avisRepo)

	// Configuration du routeur Chi
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api/avis", func(r chi.Router) {
		r.Get("/", avisHandler.GetAll)
		r.Post("/", avisHandler.Create)
		r.Get("/{id}", avisHandler.GetByID)
		r.Put("/{id}", avisHandler.Update)
		r.Delete("/{id}", avisHandler.Delete)
		r.Post("/multiple", avisHandler.CreateMultiple) 
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Serveur démarré sur le port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}