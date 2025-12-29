package auth

import (
	"log"

	"github.com/NurulloMahmud/habits/config"
	"github.com/NurulloMahmud/habits/internal/user"
)

type JWTService struct {
	cfg      config.Config
	userRepo user.Repository
	logger   *log.Logger
}

func NewJWTService(cfg config.Config, repo user.Repository, logger *log.Logger) JWTService {
	return JWTService{
		cfg:      cfg,
		userRepo: repo,
		logger:   logger,
	}
}
