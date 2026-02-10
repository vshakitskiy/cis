package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vshakitskiy/cis/internal/model"
	"github.com/vshakitskiy/cis/internal/repository"
	"github.com/vshakitskiy/cis/internal/response"
	"github.com/vshakitskiy/cis/internal/service"
)

func Auth(authService *service.AuthService, userRepo *repository.UserRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				response.WriteJSON(w, http.StatusUnauthorized, response.JSON{"error": "missing authorization header"})
				return
			}

			token := strings.TrimPrefix(header, "Bearer ")
			if token == header {
				response.WriteJSON(w, http.StatusUnauthorized, response.JSON{"error": "invalid authorization format"})
				return
			}

			userID, err := authService.ParseToken(token)
			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, response.JSON{"error": "invalid token"})
				return
			}

			user, err := userRepo.GetByID(r.Context(), userID)
			if err != nil {
				response.WriteJSON(w, http.StatusUnauthorized, response.JSON{"error": "user not found"})
				return
			}

			ctx := context.WithValue(r.Context(), model.UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
