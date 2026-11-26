package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/Gironde_API_G5/database"
	"example.com/Gironde_API_G5/handlers"
	"example.com/Gironde_API_G5/repositories"
	"example.com/Gironde_API_G5/routes"
	"example.com/Gironde_API_G5/services"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Impossible de charger le fichier .env")
	}

	client, err := database.ConnectMongoDB()

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	dbName := os.Getenv("MONGODB_DATABASE")

	if dbName == "" {
		log.Fatal("MONGODB_DATABASE manquant")
	}

	db := client.Database(dbName)

	userRepository := repositories.NewUserRepository(db)

	authService := services.NewAuthService(userRepository)

	authHandler := handlers.NewAuthHTTPHandler(authService)

	mux := http.NewServeMux()

	routes.RegisterAuthRoutes(mux, authHandler)

	mux.HandleFunc("GET /health", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("API Agence opérationnelle"))
	})

	port := os.Getenv("API_PORT")

	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:              "127.0.0.1:" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("API REST démarrée sur le port " + port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
