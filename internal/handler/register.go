package handler

import (
	"PetProject/internal/service"
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("JSON Decode Error"))
		return
	}

	err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp := RegisterResponse{
		Message: "Success",
	}
	json.NewEncoder(w).Encode(resp)
}
