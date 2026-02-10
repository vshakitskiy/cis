package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/register", h.register)
	r.Post("/login", h.login)

	return r
}

type registerRequest struct {
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Name     string         `json:"name"`
	Role     model.UserRole `json:"role"`
}

func (h *AuthHandler) register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, JSON{"error": "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		WriteJSON(w, http.StatusBadRequest, JSON{"error": "email, password, and name are required"})
		return
	}

	if req.Role == "" {
		req.Role = model.RoleEmployee
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name, req.Role)
	if errors.Is(err, service.ErrEmailTaken) {
		WriteJSON(w, http.StatusConflict, JSON{"error": err.Error()})
		return
	}
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, JSON{"error": "registration failed"})
		return
	}

	WriteJSON(w, http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteJSON(w, http.StatusBadRequest, JSON{"error": "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		WriteJSON(w, http.StatusBadRequest, JSON{"error": "email and password are required"})
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		WriteJSON(w, http.StatusUnauthorized, JSON{"error": err.Error()})
		return
	}
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, JSON{"error": "login failed"})
		return
	}

	WriteJSON(w, http.StatusOK, JSON{"token": token})
}
