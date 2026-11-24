package handlers

import (
	"context"

	"example.com/Gironde_API_G5/pb"
	"example.com/Gironde_API_G5/services"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService *services.AuthService
}

func NewAuthHandler(
	authService *services.AuthService,
) *AuthHandler {

	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(
	ctx context.Context,
	req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {

	user, err := h.authService.Register(
		ctx,
		req.GetFirstname(),
		req.GetLastname(),
		req.GetEmail(),
		req.GetPassword(),
	)

	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Inscription impossible",
		}, nil
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "Inscription réussie",
		UserId:  user.ID.Hex(),
	}, nil
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *pb.LoginRequest,
) (*pb.LoginResponse, error) {

	user, err := h.authService.Login(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
	)

	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Connexion impossible",
		}, nil
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "Connexion réussie",
		UserId:  user.ID.Hex(),
	}, nil
}
