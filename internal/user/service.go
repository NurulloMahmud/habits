package user

import (
	"context"
	"errors"
)

var (
	errEmailTaken = errors.New("This email already exists")
)

type UserService struct {
	repo Repository
}

func NewService(repo Repository) UserService {
	return UserService{
		repo: repo,
	}
}

func (s *UserService) Register(ctx context.Context, req registerUserRequest) (*User, error) {
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
