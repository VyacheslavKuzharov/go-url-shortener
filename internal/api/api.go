package api

import (
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/storage/inmemory"
	"net/http"
)

type API struct {
	router *http.ServeMux
}

func New(router *http.ServeMux) *API {
	api := API{
		router: router,
	}
	return &api
}

func (api *API) InitRoutes(storage *inmemory.MemStorage) {
	api.router.HandleFunc(`/`, saveURLHandler(storage))
	api.router.HandleFunc(`/{shortKey}`, redirectHandler(storage))
}
