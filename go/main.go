package main

import (
	"context"
	"log"

	"example.com/Gironde_API_G5/database"
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

}
