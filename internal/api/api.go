package api

import (
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api/middlewares"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"net"
)

type API struct {
	router  *chi.Mux
	cfg     *config.Config
	storage storage.Storager
	logger  *logger.Logger
}

func New(r *chi.Mux, cfg *config.Config, s storage.Storager, l *logger.Logger) *API {
	api := &API{
		router:  r,
		cfg:     cfg,
		storage: s,
		logger:  l,
	}
	api.start()

	return api
}

func (api *API) start() {
	api.router.Use(middlewares.Logger(api.logger))

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
