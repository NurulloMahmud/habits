package server

import (
	"database/sql"
	"log"
	"os"

	"github.com/NurulloMahmud/habits/config"
	"github.com/NurulloMahmud/habits/internal/platform/database"
	"github.com/NurulloMahmud/habits/internal/user"
	"github.com/NurulloMahmud/habits/migrations"
)

type Application struct {
	Logger      *log.Logger
	userHandler user.UserHandler
	DB          *sql.DB
	Cfg         config.Config
}

func NewApplication(cfg config.Config) (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	
	pgDB, err := database.New(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	err = database.Migrate(pgDB, migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	// set up users
	userRepo := user.NewRepository(pgDB)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService, logger)

	app := &Application{
		Logger:      logger,
		userHandler: *userHandler,
	}

	return app, nil
}
