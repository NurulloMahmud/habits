package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/NurulloMahmud/habits/internal/user"
	"github.com/NurulloMahmud/habits/pkg/response"
)

var (
	errUserNotFound = errors.New("User with this id not found")
)

func (s *JWTService) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, user.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.Unauthorized(w, r, "invalid authorization header")
			return
		}

		token := headerParts[1]
		claims, err := s.VerifyToken(token)
		if err != nil {
			response.Unauthorized(w, r, "invalid token")
			return
		}

		user, err := s.userRepo.Get(r.Context(), claims.ID, "")
		if err != nil {
			response.InternalServerError(w, r, err, s.logger)
			return
		}

		if user == nil {
			response.BadRequest(w, r, errUserNotFound, s.logger)
			return
		}

		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
