package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/response"
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
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "email, password, and name are required"})
		return
	}

	if req.Role == "" {
		req.Role = model.RoleEmployee
	}

	user, err := h.authService.Register(r.Context(), req.Email, req.Password, req.Name, req.Role)
	if errors.Is(err, service.ErrEmailTaken) {
		response.WriteJSON(w, http.StatusConflict, response.JSON{"error": err.Error()})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "registration failed"})
		return
	}

	response.WriteJSON(w, http.StatusCreated, user)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := response.ReadJSON(r, &req); err != nil {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "invalid request body"})
		return
	}

	if req.Email == "" || req.Password == "" {
		response.WriteJSON(w, http.StatusBadRequest, response.JSON{"error": "email and password are required"})
		return
	}

	token, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidCredentials) {
		response.WriteJSON(w, http.StatusUnauthorized, response.JSON{"error": err.Error()})
		return
	}
	if err != nil {
		response.WriteJSON(w, http.StatusInternalServerError, response.JSON{"error": "login failed"})
		return
	}

	response.WriteJSON(w, http.StatusOK, response.JSON{"token": token})
}
