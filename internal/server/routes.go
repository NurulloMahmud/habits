package server

import "github.com/go-chi/chi/v5"

func (app *Application) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/v1/register", app.userHandler.Register)

	return r
}