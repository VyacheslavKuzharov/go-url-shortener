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
	router  *chi.Mux
	cfg     *config.Config
	storage storage.Storager
}

func New(router *chi.Mux, cfg *config.Config, storage storage.Storager) *API {
	api := &API{
		router:  router,
		cfg:     cfg,
		storage: storage,
	}
	api.start()

	return api
}

func (api *API) start() {
	api.router.Use(middleware.Logger)

	api.router.Post(`/`, saveURLHandler(api.storage, api.cfg))
	api.router.Get(`/{shortKey}`, redirectHandler(api.storage))
}

func FullShortenedURL(shortKey string, cfg *config.Config) string {
	schema := "http"

	if cfg.BaseURL.Addr != "" {
		return fmt.Sprintf("%s/%s", cfg.BaseURL.Addr, shortKey)
	}

	addr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	return fmt.Sprintf("%s://%s/%s", schema, addr, shortKey)
}
