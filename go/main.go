package main

import (
	"context"
	"log"
	"net"
	"os"

	"example.com/Gironde_API_G5/database"
	"example.com/Gironde_API_G5/handlers"
	"example.com/Gironde_API_G5/pb"
	"example.com/Gironde_API_G5/repositories"
	"example.com/Gironde_API_G5/services"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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

	authHandler := handlers.NewAuthHandler(authService)

	listener, err := net.Listen("tcp", "127.0.0.1:50051")

	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	log.Println("Serveur gRPC démarré sur le port 50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
