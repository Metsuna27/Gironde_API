package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"example.com/Gironde_API_G5/dto"
	"example.com/Gironde_API_G5/utils"
)

type contextKey string

const claimsKey contextKey = "auth_claims"

func AuthMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {

		authHeader := r.Header.Get("Authorization")

		parts := strings.Fields(authHeader)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") {

			writeError(
				w,
				http.StatusUnauthorized,
				"Token manquant ou format invalide",
			)

			return
		}

		tokenString := parts[1]

		claims, err := utils.VerifyToken(tokenString)

		if err != nil {

			writeError(
				w,
				http.StatusUnauthorized,
				"Token invalide ou expiré",
			)

			return
		}

		ctx := context.WithValue(
			r.Context(),
			claimsKey,
			claims,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(ctx context.Context) (*utils.JWTClaims, bool) {

	claims, ok := ctx.Value(claimsKey).(*utils.JWTClaims)

	return claims, ok
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(dto.AuthResponse{
		Success: false,
		Message: message,
	})
}
