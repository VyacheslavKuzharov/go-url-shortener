package api

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	router *chi.Mux
}

func New(router *chi.Mux) *API {
	api := API{
		router: router,
	}
	return &api
}

func (api *API) InitRoutes(storage storage.Storager) chi.Router {
	api.router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Post(`/`, saveURLHandler(storage))
		r.Get(`/{shortKey}`, redirectHandler(storage))
	})

	return api.router
}
