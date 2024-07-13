package config

import (
	"errors"
	"strings"
)

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
