package utils

import (
	"errors"
	"os"
	"time"

	"example.com/Gironde_API_G5/models"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateToken(user *models.User) (string, error) {

	secret := os.Getenv("JWT_SECRET")

	if len(secret) < 32 {
		return "", errors.New("JWT_SECRET invalide")
	}

	if user.Role != models.RoleAdmin &&
		user.Role != models.RoleVendeur {
		return "", errors.New("rôle utilisateur invalide")
	}

	now := time.Now()

	claims := JWTClaims{
		Role: user.Role,

		RegisteredClaims: jwt.RegisteredClaims{
			Subject:  user.ID.Hex(),
			Issuer:   "richnou-immo-agence",
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

func VerifyToken(tokenString string) (*JWTClaims, error) {

	secret := os.Getenv("JWT_SECRET")

	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET invalide")
	}

	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		},
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuer("richnou-immo-agence"),
	)

	if err != nil || !token.Valid {
		return nil, errors.New("token invalide ou expiré")
	}

	if claims.Subject == "" {
		return nil, errors.New("identifiant utilisateur manquant")
	}

	if claims.Role != models.RoleAdmin &&
		claims.Role != models.RoleVendeur {

		return nil, errors.New("rôle invalide")
	}

	return claims, nil
}
