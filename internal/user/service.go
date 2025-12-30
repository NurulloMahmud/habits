package user

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/NurulloMahmud/habits/config"
	"github.com/NurulloMahmud/habits/internal/auth"
)

var (
	errEmailTaken         = errors.New("This email already exists")
	errInvalidCredentials = errors.New("invalid credentials")
	errUserInactive       = errors.New("inactive user")
	errUserLocked         = errors.New("user is locked temporarily")
	errMatchingPassword   = errors.New("error matching password")
)

type UserService struct {
	repo Repository
	cfg  config.Config
}

func NewService(repo Repository, cfg config.Config) UserService {
	return UserService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *UserService) register(ctx context.Context, req registerUserRequest) (*User, error) {
	existingUser, err := s.repo.Get(ctx, 0, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errEmailTaken
	}

	newUser := User{Email: req.Email, UserRole: "user"}
	if req.FirstName != nil {
		newUser.FirstName = req.FirstName
	}
	if req.LastName != nil {
		newUser.LastName = req.LastName
	}

	err = newUser.PasswordHash.Set(req.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) login(ctx context.Context, email, password string) (*User, string, error) {
	user, err := s.repo.Get(ctx, 0, email)
	if err != nil {
		return nil, "", err
	}

	if user == nil {
		return nil, "", errInvalidCredentials
	}

	if !user.IsActive {
		return nil, "", errUserInactive
	}

	if user.IsLocked {
		if time.Since(user.LastFailedLogin.Time) < time.Hour*24 {
			return nil, "", errUserLocked
		}
		if err := s.repo.Unlock(ctx, user.ID); err != nil {
			return nil, "", err
		}
	}

	matched, err := user.PasswordHash.Matches(password)
	if err != nil {
		return nil, "", errMatchingPassword
	}

	if !matched {
		user.FailedAttempts += 1
		user.LastFailedLogin = sql.NullTime{
			Time:  time.Now().UTC(),
			Valid: true,
		}

		if user.FailedAttempts >= 5 {
			user.IsLocked = true
		}

		err = s.repo.Update(ctx, *user)
		if err != nil {
			return nil, "", err
		}

		return nil, "", errInvalidCredentials
	}

	user.FailedAttempts = 0
	err = s.repo.Update(ctx, *user)
	if err != nil {
		return nil, "", err
	}

	claims := auth.TokenClaims{
		ID:       user.ID,
		Email:    user.Email,
		UserRole: user.UserRole,
	}

	token, err := auth.GenerateAccessToken(claims, s.cfg.JWTSecret)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) update(ctx context.Context, id int64, req updateUserRequest) error {
	user, err := s.repo.Get(ctx, id, "")
	if err != nil {
		return err
	}

	matched, err := user.PasswordHash.Matches(*req.OldPassword)
	if err != nil {
		return err
	}

	if !matched {
		return errInvalidCredentials
	}

	if req.Email != nil {
		existingUser, err := s.repo.Get(ctx, 0, *req.Email)
		if err != nil {
			return err
		}

		if existingUser != nil && existingUser.ID != user.ID {
			return errEmailTaken
		}
		user.Email = *req.Email
	}

	if req.FirstName != nil {
		user.FirstName = req.FirstName
	}
	if req.LastName != nil {
		user.LastName = req.LastName
	}
	if req.NewPassword != nil {
		err = user.PasswordHash.Set(*req.NewPassword)
		if err != nil {
			return err
		}
	}

	err = s.repo.Update(ctx, *user)
	if err != nil {
		return err
	}

	return nil
}
