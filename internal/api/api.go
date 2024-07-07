package api

import (
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net"
)

type API struct {
	router *chi.Mux
	cfg    *config.Config
}

func New(router *chi.Mux, cfg *config.Config) *API {
	api := API{
		router: router,
		cfg:    cfg,
	}
	return &api
}

func (api *API) InitRoutes(storage storage.Storager) chi.Router {
	api.router.Group(func(r chi.Router) {
		r.Use(middleware.Logger)

		r.Post(`/`, saveURLHandler(storage, api.cfg))
		r.Get(`/{shortKey}`, redirectHandler(storage))
	})

	return api.router
}

func FullShortenedURL(shortKey string, cfg *config.Config) string {
	schema := "http"

	if cfg.BaseURL.Addr != "" {
		return fmt.Sprintf("%s/%s", cfg.BaseURL.Addr, shortKey)
	}

	addr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	return fmt.Sprintf("%s://%s/%s", schema, addr, shortKey)
}
