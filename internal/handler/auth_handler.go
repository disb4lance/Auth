package handler

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/service"
)

type AuthHandler struct {
	authService service.AuthService
}

// Конструктор
func NewAuthHandler(s service.AuthService) *AuthHandler {
	return &AuthHandler{authService: s}
}

// RegisterHandler godoc
// @Summary Register a new user
// @Description Create a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "User info"
// @Success 201 {object} service.UserDTO
// @Failure 400 {string} string "invalid body"
// @Router /auth/register [post]
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
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
	json.NewEncoder(w).Encode(userDTO)
}

// AuthenticateHandler godoc
// @Summary Authenticate user
// @Description Login user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body service.Credentials true "User credentials"
// @Success 200 {object} service.AuthenticatedUser
// @Failure 401 {string} string "unauthorized"
// @Router /auth/tokens [post]
func (h *AuthHandler) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds service.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	authUser, err := h.authService.Authenticate(creds)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authUser)
}

// RefreshHandler godoc
// @Summary Refresh JWT token
// @Description Refresh access and refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body map[string]string true "Refresh token"
// @Success 200 {object} service.AuthenticatedUser
// @Failure 401 {string} string "invalid refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(authUser)
}
