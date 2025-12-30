package user

import (
	"errors"
	"regexp"
	"strings"
)

var (
	errEmailFormat                = errors.New("invalid email")
	errEmailRequired              = errors.New("email is required")
	errPasswordRequired           = errors.New("password is required")
	errPasswordAndConfirmRequired = errors.New("password and password_confirm fields are required")
	errPasswordLen                = errors.New("Password must be between 6 and 32 characters long")
	errPasswordsNotMatch          = errors.New("Passwords do not match")
	errFirstNameEmpty             = errors.New("you can omit first_name but cannot send empty string or space")
	errLastNameEmpty              = errors.New("you can omit last_name but cannot send empty string or space")
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
		return errPasswordAndConfirmRequired
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

type loginRequest struct {
	Email    *string `jsoon:"email"`
	Password *string `json:"password"`
}

func (r *loginRequest) validateLoginRequest() error {
	if r.Email == nil {
		return errEmailRequired
	} else {
		err := validateEmailFormat(*r.Email)
		if err != nil {
			return errEmailFormat
		}
	}

	if r.Password == nil {
		return errPasswordRequired
	}

	return nil
}

type updateUserRequest struct {
	Email              *string `json:"email"`
	FirstName          *string `json:"first_name"`
	LastName           *string `json:"last_name"`
	OldPassword        *string `json:"old_password"`
	NewPassword        *string `json:"new_password"`
	NewPasswordConfirm *string `json:"new_password_confirm"`
}

func (r *updateUserRequest) validateUpdateUserRequest() error {
	if r.Email != nil {
		err := validateEmailFormat(*r.Email)
		if err != nil {
			return errEmailFormat
		}
	}
	return nil
}

func (r *updateUserRequest) validatePasswordUpdate() error {
	if r.OldPassword == nil {
		return errPasswordRequired
	}
	if r.NewPassword == nil || r.NewPasswordConfirm == nil {
		return errors.New("new_password and new_password_confirm fields are required")
	}
	if len(*r.NewPassword) > 32 || len(*r.NewPassword) < 6 {
		return errPasswordLen
	}
	if r.NewPassword != r.NewPasswordConfirm {
		return errPasswordsNotMatch
	}
	return nil
}
