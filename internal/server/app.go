package server

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/NurulloMahmud/habits/config"
	"github.com/NurulloMahmud/habits/internal/habit"
	habitmember "github.com/NurulloMahmud/habits/internal/habit_member"
	"github.com/NurulloMahmud/habits/internal/middleware"
	"github.com/NurulloMahmud/habits/internal/platform/database"
	"github.com/NurulloMahmud/habits/internal/user"
	"github.com/NurulloMahmud/habits/migrations"
	"github.com/NurulloMahmud/habits/pkg/context"
	"github.com/NurulloMahmud/habits/pkg/response"
)

type Application struct {
	Logger             *log.Logger
	userHandler        user.UserHandler
	habitHandler       habit.HabitHandler
	habitMemberHandler habitmember.Handler
	DB                 *sql.DB
	Cfg                config.Config
	middleware         middleware.Middleware
}

func NewApplication(cfg config.Config) (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// spin up postgres db
	pgDB, err := database.New(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// spin up mongodb
	err = database.ConnectMongo(cfg.MongoDBURL)
	if err != nil {
		log.Fatal(err)
	}

	err = database.Migrate(pgDB, migrations.FS, ".")
	if err != nil {
		return nil, err
	}

	// set up repositories
	userRepo := user.NewPostgresRepository(pgDB)
	habitRepo := habit.NewPostgresRepository(pgDB)
	habitMemberRepo := habitmember.NewPostgresRepository(pgDB)

	// setup services
	userService := user.NewService(userRepo, cfg)
	habitService := habit.NewHabitService(habitRepo)
	habitMemberService := habitmember.NewService(habitMemberRepo)

	// setup handlers
	userHandler := user.NewHandler(userService, logger)
	habitHandler := habit.NewHandler(habitService, logger)
	habitMemberHandler := habitmember.NewHandler(habitMemberService, logger)

	// setup middlewares
	appMiddleware := middleware.NewMiddleware(logger, userRepo, cfg)

	app := &Application{
		Logger:             logger,
		userHandler:        *userHandler,
		habitHandler:       *habitHandler,
		habitMemberHandler: *habitMemberHandler,
		middleware:         *appMiddleware,
		DB:                 pgDB,
		Cfg:                cfg,
	}

	return app, nil
}

func (a *Application) testHandler(w http.ResponseWriter, r *http.Request) {
	user := context.GetUser(r)
	response.WriteJSON(w, http.StatusOK, response.Envelope{"user": user})
}

func (a *Application) health(w http.ResponseWriter, r *http.Request) {
	response.WriteJSON(w, http.StatusOK, response.Envelope{"status": "available"})
}
