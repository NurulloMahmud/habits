package context

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type contextKey string

const userContextKey = contextKey("user")

type User struct {
	ID              int64        `json:"id"`
	Email           string       `json:"email"`
	FirstName       *string      `json:"first_name"`
	LastName        *string      `json:"last_name"`
	UserRole        string       `json:"user_role"`
	IsActive        bool         `json:"is_active"`
	IsLocked        bool         `json:"is_locked"`
	LastFailedLogin sql.NullTime `json:"-"`
	FailedAttempts  int64        `json:"-"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func SetUser(r *http.Request, user *User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *User {
	user, ok := r.Context().Value(userContextKey).(*User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
