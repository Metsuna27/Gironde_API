package repositories

import (
	"context"
	"fmt"

	"example.com/Gironde_API_G5/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	user *models.User,
) error {

	fmt.Printf("Email avant insertion : %s\n", user.Email)
	fmt.Printf("Role avant insertion : %q\n", user.Role)

	_, err := r.collection.InsertOne(ctx, user)

	return err
}

func (r *UserRepository) FindUserByEmail(
	ctx context.Context,
	email string,
) (*models.User, error) {

	var user models.User

	err := r.collection.FindOne(
		ctx,
		bson.M{"email": email},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindUserByID(
	ctx context.Context,
	id bson.ObjectID,
) (*models.User, error) {

	var user models.User

	err := r.collection.FindOne(
		ctx,
		bson.M{"_id": id},
	).Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
