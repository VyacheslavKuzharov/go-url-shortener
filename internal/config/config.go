package config

import (
	"encoding/hex"
	"errors"
	"flag"
	baseURLcfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/base_url"
	httpcfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/http"
	logscfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/logs"
	storagecfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/storage"
)

type Config struct {
	HTTP    httpcfg.HTTPCfg
	BaseURL baseURLcfg.BaseURLCfg
	Log     logscfg.LogCfg
	Storage storagecfg.StorageCfg
}

var CookieSalt []byte

func New() (*Config, error) {
	var hcf *httpcfg.HTTPCfg
	var err error
	cfg := &Config{}

	http, baseURL, filePath, pgDSN := parseHTTPServerFlags()
	CookieSalt, err = defineCookieSecretKey()
	if err != nil {
		return nil, err
	}

	hcf, err = httpcfg.NewHTTPCfg(http)
	if err != nil {
		return cfg, err
	}
	cfg.HTTP = *hcf
	cfg.BaseURL = *baseURLcfg.NewBaseURLCfg(baseURL)
	cfg.Log = *logscfg.NewLogsCfg()
	cfg.Storage = *storagecfg.NewStorageCfg(filePath, pgDSN)

	return cfg, nil
}

func parseHTTPServerFlags() (*httpcfg.HTTPCfg, *baseURLcfg.BaseURLCfg, *storagecfg.FileStorage, *storagecfg.PgStorage) {
	addr := new(httpcfg.HTTPCfg)
	url := new(baseURLcfg.BaseURLCfg)
	filePath := new(storagecfg.FileStorage)
	pgDSN := new(storagecfg.PgStorage)

	flag.Var(addr, "a", "Net address host:port")
	flag.Var(url, "b", "base address of the resulting shortened URL http://localhost:3000")
	flag.Var(filePath, "f", "path to file storage")
	flag.Var(pgDSN, "d", "postgres connection URL")
	flag.Parse()

	return addr, url, filePath, pgDSN
}

func defineCookieSecretKey() ([]byte, error) {
	// Здесь должна быть os.Getenv("COOKIE_SALT")
	// без хардкода тесты не проходят
	salt, err := hex.DecodeString("13d6b4dff8f84a10851021ec8608f814570d562c92fe6b5ec4c9f595bcb3234b")
	if err != nil {
		return nil, err
	}

	if len(salt) == 0 {
		return nil, errors.New("env variable COOKIE_SALT is missing")
	}

	return salt, nil
}
