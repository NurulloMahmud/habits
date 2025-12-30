package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/NurulloMahmud/habits/config"
	"github.com/NurulloMahmud/habits/internal/middleware"
	"github.com/NurulloMahmud/habits/internal/platform/database"
	"github.com/NurulloMahmud/habits/internal/user"
	"github.com/NurulloMahmud/habits/migrations"
	"github.com/NurulloMahmud/habits/pkg/response"
)

type Application struct {
	Logger      *log.Logger
	userHandler user.UserHandler
	DB          *sql.DB
	Cfg         config.Config
	middleware  middleware.Middleware
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

	// set up repositories
	userRepo := user.NewRepository(pgDB)

	// setup services
	userService := user.NewService(userRepo, cfg)

	// setup handlers
	userHandler := user.NewHandler(userService, logger)

	// setup middlewares
	appMiddleware := middleware.NewMiddleware(logger, userRepo, cfg)

	app := &Application{
		Logger:      logger,
		userHandler: *userHandler,
		middleware:  *appMiddleware,
		DB:          pgDB,
		Cfg:         cfg,
	}

	return app, nil
}

func (a *Application) testHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)
	a.Logger.Printf("is anonymous: %t", user.IsAnonymous())
	response.WriteJSON(w, http.StatusOK, response.Envelope{"user": user})
}
