package middleware

import (
	"log"

	"github.com/NurulloMahmud/habits/config"
	"github.com/NurulloMahmud/habits/internal/user"
)

type Middleware struct {
	logger   *log.Logger
	userRepo user.Repository
	cfg      config.Config
}

func NewMiddleware(logger *log.Logger, repo user.Repository, cfg config.Config) *Middleware {
	return &Middleware{
		logger:   logger,
		userRepo: repo,
		cfg:      cfg,
	}
}
