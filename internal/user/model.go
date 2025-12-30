package user

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

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
	PasswordHash    password     `json:"-"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type password struct {
	plaintText *string
	hash       []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintText = &plaintextPassword
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err //internal server error
		}
	}

	return true, nil
}
