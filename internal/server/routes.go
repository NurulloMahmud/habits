package server

import "github.com/go-chi/chi/v5"

func (app *Application) Routes() *chi.Mux {
	r := chi.NewRouter()

	// register & login
	r.Post("/api/v1/register", app.userHandler.Register)
	r.Post("/api/v1/login", app.userHandler.Login)

	r.Group(func(r chi.Router) {
		r.Use(app.middleware.Authenticate)
		r.Get("/test", app.testHandler)

		r.Use(app.middleware.RequireUser)
		r.Patch("/api/v1/update", app.userHandler.Update)
	})

	return r
}
