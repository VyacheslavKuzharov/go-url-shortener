package config

import (
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/httpserver"
	"os"
	"strings"
)

type HTTPCfg struct {
	Host string
	Port string
}

func NewHTTPCfg(flagHTTP *HTTPCfg) *HTTPCfg {
	if flagHTTP.Host != "" && flagHTTP.Port != "" {
		return flagHTTP
	}

	cfg := &HTTPCfg{
		Host: httpserver.DefaultHost,
		Port: httpserver.DefaultPort,
	}

	if os.Getenv("SERVER_ADDRESS") != "" {
		hp := strings.Split(os.Getenv("SERVER_ADDRESS"), ":")

		cfg.Host = hp[0]
		cfg.Port = hp[1]
	}

	return cfg
}

func (a *HTTPCfg) String() string {
	return a.Host + ":" + a.Port
}

func (a *HTTPCfg) Set(s string) error {
	hp := strings.Split(s, ":")
	if len(hp) != 2 {
		return errors.New("need address in a form host:port")
	}

	if hp[0] == "" {
		return errors.New("http Host can't be blank")
	}

	if hp[1] == "" {
		return errors.New("http Port can't be blank")
	}

	a.Host = hp[0]
	a.Port = hp[1]
	return nil
}
