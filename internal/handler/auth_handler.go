package handler

import (
	"auth-service/internal/handler/dto"
	"auth-service/internal/service"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Create a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User info"
// @Success 201 {object} service.UserDTO
// @Failure 400 {string} string "invalid body"
// @Router /auth/register [post]
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	userDTO, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(userDTO)
}

// AuthenticateHandler godoc
// @Summary Authenticate user
// @Description Login user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "User credentials"
// @Success 200 {object} service.AuthenticatedUser
// @Failure 401 {string} string "unauthorized"
// @Router /auth/tokens [post]
func (h *AuthHandler) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	authUser, err := h.authService.Authenticate(service.Credentials{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authUser)
}

// RefreshHandler godoc
// @Summary Refresh JWT token
// @Description Refresh access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} service.AuthenticatedUser
// @Failure 401 {string} string "invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.RefreshRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	authUser, err := h.authService.Refresh(req.RefreshToken)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authUser)
}
