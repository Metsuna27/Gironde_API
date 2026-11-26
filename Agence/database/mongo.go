package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB() (*mongo.Client, error) {

	uri := os.Getenv("MONGODB_URI")

	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI manquant")
	}

	client, err := mongo.Connect(
		options.Client().ApplyURI(uri),
	)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		client.Disconnect(context.Background())
		return nil, err
	}

	fmt.Println("Connexion à MongoDB réussie !")

	return client, nil
}
