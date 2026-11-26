package routes

import (
	"encoding/json"
	"net/http"

	"example.com/Gironde_API_G5/handlers"
	"example.com/Gironde_API_G5/middleware"
	"example.com/Gironde_API_G5/models"
)

func RegisterAuthRoutes(
	mux *http.ServeMux,
	authHandler *handlers.AuthHTTPHandler,
) {

	mux.HandleFunc(
		"POST /api/auth/register",
		authHandler.Register,
	)

	mux.HandleFunc(
		"POST /api/auth/login",
		authHandler.Login,
	)

	mux.Handle(
		"GET /api/auth/me",
		middleware.AuthMiddleware(
			http.HandlerFunc(meHandler),
		),
	)

	mux.Handle(
		"GET /api/admin/test",
		middleware.AuthMiddleware(
			middleware.RoleMiddleware(models.RoleAdmin)(
				http.HandlerFunc(adminHandler),
			),
		),
	)

}

func meHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	claims, _ := middleware.GetClaims(r.Context())

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Utilisateur authentifié",
		"user_id": claims.Subject,
		"role":    claims.Role,
	})
}

func adminHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "Bienvenue dans l'espace administrateur",
	})
}
