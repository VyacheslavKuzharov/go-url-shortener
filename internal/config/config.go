package config

import (
	"flag"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/httpserver"
	"os"
	"strings"
)

type Config struct {
	HTTP    HTTPCfg
	BaseURL BaseURLCfg
	Log     LogCfg
}

type HTTPCfg struct {
	Host string
	Port string
}

type BaseURLCfg struct {
	Addr string
}

type LogCfg struct {
	Level string
}

func New() (*Config, error) {
	cfg := &Config{}
	httpcfg, baseURL := parseHTTPServerFlags()

	if os.Getenv("SERVER_ADDRESS") != "" {
		hp := strings.Split(os.Getenv("SERVER_ADDRESS"), ":")

		cfg.HTTP.Host = hp[0]
		cfg.HTTP.Port = hp[1]
	} else if httpcfg.Host != "" && httpcfg.Port != "" {
		cfg.HTTP = *httpcfg
	} else {
		cfg.HTTP.Host = httpserver.DefaultHost
		cfg.HTTP.Port = httpserver.DefaultPort
	}

	if os.Getenv("BASE_URL") != "" {
		cfg.BaseURL.Addr = os.Getenv("BASE_URL")
	} else if baseURL.Addr != "" {
		cfg.BaseURL = *baseURL
	}

	cfg.Log.Level = "info"

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
