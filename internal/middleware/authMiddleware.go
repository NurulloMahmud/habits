package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/NurulloMahmud/habits/internal/auth"
	"github.com/NurulloMahmud/habits/pkg/context"
	"github.com/NurulloMahmud/habits/pkg/response"
)

var (
	errUserNotFound = errors.New("user not found")
)

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = context.SetUser(r, context.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.Unauthorized(w, r, "invalid authorization header")
			return
		}

		token := headerParts[1]
		claims, err := auth.VerifyToken(token, m.cfg.JWTSecret)
		if err != nil {
			response.Unauthorized(w, r, "invalid token")
			return
		}

		user, err := m.userRepo.Get(r.Context(), claims.ID, "")
		if err != nil {
			response.InternalServerError(w, r, err, m.logger)
			return
		}

		if user == nil {
			response.BadRequest(w, r, errUserNotFound, m.logger)
			return
		}

		contextUser := context.User{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			UserRole:  user.UserRole,
			IsActive:  user.IsActive,
			IsLocked:  user.IsLocked,
		}

		r = context.SetUser(r, &contextUser)
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userContext := context.GetUser(r)
		if userContext.IsAnonymous() || !userContext.IsActive || userContext.IsLocked {
			response.Unauthorized(w, r, "Unauthoized")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAdminUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userContext := context.GetUser(r)
		if userContext.IsAnonymous() || !userContext.IsActive || userContext.IsLocked {
			response.Unauthorized(w, r, "Unauthoized")
			return
		}

		user, err := m.userRepo.Get(r.Context(), userContext.ID, "")
		if err != nil {
			response.InternalServerError(w, r, err, m.logger)
			return
		}

		if user.UserRole != "admin" {
			response.Unauthorized(w, r, "Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	})
}
