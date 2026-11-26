package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/Gironde_API_G5/dto"
	"example.com/Gironde_API_G5/services"
	"example.com/Gironde_API_G5/utils"
)

type AuthHTTPHandler struct {
	authService *services.AuthService
}

func NewAuthHTTPHandler(
	authService *services.AuthService,
) *AuthHTTPHandler {

	return &AuthHTTPHandler{
		authService: authService,
	}
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}

func (h *AuthHTTPHandler) Register(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.RegisterRequest

	r.Body = http.MaxBytesReader(w, r.Body, 16384)

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {

		writeJSON(w, http.StatusBadRequest,
			dto.AuthResponse{
				Success: false,
				Message: "Requête JSON invalide",
			},
		)

		return
	}

	user, err := h.authService.Register(
		r.Context(),
		req.Firstname,
		req.Lastname,
		req.Email,
		req.Password,
	)

	if err != nil {

		writeJSON(w, http.StatusBadRequest,
			dto.AuthResponse{
				Success: false,
				Message: "Inscription impossible",
			},
		)

		return
	}

	writeJSON(w, http.StatusCreated,
		dto.AuthResponse{
			Success: true,
			Message: "Inscription réussie",
			UserID:  user.ID.Hex(),
			Role:    user.Role,
		},
	)
}

func (h *AuthHTTPHandler) Login(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.LoginRequest

	r.Body = http.MaxBytesReader(w, r.Body, 16384)

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		writeJSON(w, http.StatusBadRequest,
			dto.AuthResponse{
				Success: false,
				Message: "Requête JSON invalide",
			},
		)

		return
	}

	// Vérification des identifiants
	user, err := h.authService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)

	if err != nil {
		writeJSON(w, http.StatusUnauthorized,
			dto.AuthResponse{
				Success: false,
				Message: "Connexion impossible",
			},
		)

		return
	}

	// Génération du token JWT
	token, err := utils.GenerateToken(user)

	if err != nil {
		writeJSON(w, http.StatusInternalServerError,
			dto.AuthResponse{
				Success: false,
				Message: "Impossible de générer le token",
			},
		)

		return
	}

	// Réponse au client
	writeJSON(w, http.StatusOK,
		dto.AuthResponse{
			Success: true,
			Message: "Connexion réussie",
			UserID:  user.ID.Hex(),
			Role:    user.Role,
			Token:   token,
		},
	)
}
