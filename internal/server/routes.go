package server

import "github.com/go-chi/chi/v5"

func (app *Application) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(app.middleware.RateLimit)

	// test & health
	r.Get("/health", app.health)

	r.Group(func(r chi.Router) {
		r.Use(app.middleware.Authenticate)
		r.Use(app.middleware.ActivityLogger)
		r.Get("/test", app.userHandler.List)

		// register & login
		r.Post("/api/v1/register", app.userHandler.Register)
		r.Post("/api/v1/login", app.userHandler.Login)

		// habits (public)
		r.Get("/api/v1/habits", app.habitHandler.HandleGetHabitList)

		// valid user required endpoints
		r.Group(func(r chi.Router) {
			r.Use(app.middleware.RequireUser)

			// users endpoints
			r.Patch("/api/v1/users", app.userHandler.Update)

			// habits
			r.Post("/api/v1/habits", app.habitHandler.HandleCreate)
			r.Patch("/api/v1/habits/{id}", app.habitHandler.HandleUpdate)
			r.Delete("/api/v1/habits/{id}", app.habitHandler.HandleDelete)

			// habit members endpoints
			r.Post("/api/v1/join-habit", app.habitMemberHandler.HandleJoinHabit)
		})
	})

	return r
}
