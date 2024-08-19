package api

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/api/middlewares"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
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
	api.router.Use(middlewares.Compress)

	api.router.Get(`/ping`, pingHandler(api.storage))
	api.router.Get(`/{shortKey}`, redirectHandler(api.storage))
	api.router.
		With(middlewares.Cookies(api.logger)).
		Post(`/`, saveURLHandler(api.storage, api.cfg))

	api.router.Route("/api", func(route chi.Router) {
		route.Group(func(public chi.Router) {
			public.Use(middlewares.Cookies(api.logger))

			public.Post(`/shorten`, shortenHandler(api.storage, api.cfg))
			public.Post(`/shorten/batch`, batchHandler(api.storage, api.cfg))
		})

		route.Group(func(private chi.Router) {
			private.Use(middlewares.Auth)

			private.Get("/user/urls", userURLsHandler(api.storage, api.cfg))
		})
	})
}
