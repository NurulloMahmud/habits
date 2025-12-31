package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/NurulloMahmud/habits/config"
	"github.com/NurulloMahmud/habits/internal/habit"
	"github.com/NurulloMahmud/habits/internal/middleware"
	"github.com/NurulloMahmud/habits/internal/platform/database"
	"github.com/NurulloMahmud/habits/internal/user"
	"github.com/NurulloMahmud/habits/migrations"
	"github.com/NurulloMahmud/habits/pkg/context"
	"github.com/NurulloMahmud/habits/pkg/response"
)

type Application struct {
	Logger       *log.Logger
	userHandler  user.UserHandler
	habitHandler habit.HabitHandler
	DB           *sql.DB
	Cfg          config.Config
	middleware   middleware.Middleware
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
	habitRepo := habit.NewRepo(pgDB)

	// setup services
	userService := user.NewService(userRepo, cfg)
	habitService := habit.NewHabitService(habitRepo)

	// setup handlers
	userHandler := user.NewHandler(userService, logger)
	habitHandler := habit.NewHandler(habitService, logger)

	// setup middlewares
	appMiddleware := middleware.NewMiddleware(logger, userRepo, cfg)

	app := &Application{
		Logger:       logger,
		userHandler:  *userHandler,
		habitHandler: *habitHandler,
		middleware:   *appMiddleware,
		DB:           pgDB,
		Cfg:          cfg,
	}

	return app, nil
}

func (a *Application) testHandler(w http.ResponseWriter, r *http.Request) {
	user := context.GetUser(r)
	response.WriteJSON(w, http.StatusOK, response.Envelope{"user": user})
}
