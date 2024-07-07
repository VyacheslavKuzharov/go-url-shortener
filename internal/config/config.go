package config

import (
	"flag"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/httpserver"
)

type Config struct {
	HTTP    HTTPCfg
	BaseURL BaseURLCfg
}

type HTTPCfg struct {
	Host string
	Port string
}

type BaseURLCfg struct {
	Addr string
}

func New() (*Config, error) {
	cfg := &Config{}
	httpcfg, baseURL := parseHTTPServerFlags()

	if httpcfg.Host != "" && httpcfg.Port != "" {
		cfg.HTTP = *httpcfg
	} else {
		cfg.HTTP.Host = httpserver.DefaultHost
		cfg.HTTP.Port = httpserver.DefaultPort
	}

	if baseURL.Addr != "" {
		cfg.BaseURL = *baseURL
	}

	return cfg, nil
}

func parseHTTPServerFlags() (*HTTPCfg, *BaseURLCfg) {
	addr := new(HTTPCfg)
	url := new(BaseURLCfg)

	flag.Var(addr, "a", "Net address host:port")
	flag.Var(url, "b", "base address of the resulting shortened URL http://localhost:3000")
	flag.Parse()

	return addr, url
}
