package httpapi

import (
	"fmt"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"net"
	"net/url"
)

func FullShortenedURL(shortKey string, cfg *config.Config) string {
	schema := "http"

	if cfg.BaseURL.Addr != "" {
		return fmt.Sprintf("%s/%s", cfg.BaseURL.Addr, shortKey)
	}

	addr := net.JoinHostPort(cfg.HTTP.Host, cfg.HTTP.Port)
	return fmt.Sprintf("%s://%s/%s", schema, addr, shortKey)
}

func IsURLValid(originalURL string) bool {
	u, err := url.Parse(originalURL)

	return err == nil && u.Scheme != "" && u.Host != ""
}
