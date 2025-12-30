package middleware

import (
	"context"
	"net/http"

	"github.com/NurulloMahmud/habits/internal/user"
)

type contextKey string

const userContextKey = contextKey("user")

func SetUser(r *http.Request, user *user.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *user.User {
	user, ok := r.Context().Value(userContextKey).(*user.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
