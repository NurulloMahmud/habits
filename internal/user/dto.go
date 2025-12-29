package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errEmailFormat       = errors.New("invalid email")
	errPasswordRequired  = errors.New("password and password_confirm fields are required")
	errPasswordLen       = errors.New("Password must be between 6 and 32 characters long")
	errPasswordsNotMatch = errors.New("Passwords do not match")
	errFirstNameEmpty    = errors.New("you can omit first_name but cannot send empty string or space")
	errLastNameEmpty     = errors.New("you can omit last_name but cannot send empty string or space")
)

type registerUserRequest struct {
	Email           string  `json:"email"`
	Password        string  `json:"password"`
	PasswordConfirm string  `json:"password_confirm"`
	FirstName       *string `json:"first_name"`
	LastName        *string `json:"last_name"`
}

func (u *registerUserRequest) validateRegister() error {
	err := validateEmailFormat(u.Email)
	if err != nil {
		return errEmailFormat
	}

	if u.Password == "" || u.PasswordConfirm == "" {
		return errPasswordRequired
	}

	if len(u.Password) < 6 || len(u.Password) > 32 {
		return errPasswordLen
	}

	if u.Password != u.PasswordConfirm {
		return errPasswordsNotMatch
	}

	if u.FirstName != nil {
		if strings.TrimSpace(*u.FirstName) == "" {
			return errFirstNameEmpty
		}
	}
	if u.LastName != nil {
		if strings.TrimSpace(*u.LastName) == "" {
			return errLastNameEmpty
		}
	}

	return nil
}

func validateEmailFormat(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errEmailFormat
	}
	return nil
}
