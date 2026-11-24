package services

import (
	"context"
	"errors"
	"net/mail"
	"strings"
	"time"

	"example.com/Gironde_API_G5/models"
	"example.com/Gironde_API_G5/repositories"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepository *repositories.UserRepository
}

func NewAuthService(
	userRepository *repositories.UserRepository,
) *AuthService {

	return &AuthService{
		userRepository: userRepository,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	firstname string,
	lastname string,
	email string,
	password string,
) (*models.User, error) {

	firstname = strings.TrimSpace(firstname)
	lastname = strings.TrimSpace(lastname)
	email = strings.ToLower(strings.TrimSpace(email))

	if firstname == "" || lastname == "" ||
		email == "" || password == "" {

		return nil, errors.New(
			"tous les champs sont obligatoires",
		)
	}

	address, err := mail.ParseAddress(email)

	if err != nil || address.Address != email {
		return nil, errors.New("adresse email invalide")
	}

	if len(password) < 8 || len(password) > 72 {
		return nil, errors.New(
			"le mot de passe doit contenir entre 8 et 72 octets",
		)
	}

	_, err = s.userRepository.FindUserByEmail(ctx, email)

	if err == nil {
		return nil, errors.New("email déjà utilisé")
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	user := &models.User{
		ID:        bson.NewObjectID(),
		Firstname: firstname,
		Lastname:  lastname,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.userRepository.CreateUser(ctx, user)

	if mongo.IsDuplicateKeyError(err) {
		return nil, errors.New("email déjà utilisé")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Login(
	ctx context.Context,
	email string,
	password string,
) (*models.User, error) {

	email = strings.ToLower(strings.TrimSpace(email))

	if email == "" || password == "" {
		return nil, errors.New(
			"email et mot de passe obligatoires",
		)
	}

	user, err := s.userRepository.FindUserByEmail(ctx, email)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errors.New("identifiants incorrects")
	}

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)

	if err != nil {
		return nil, errors.New("identifiants incorrects")
	}

	return user, nil
}
