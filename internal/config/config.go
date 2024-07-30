package config

import (
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

func New() (*Config, error) {
	cfg := &Config{}
	http, baseURL, filePath := parseHTTPServerFlags()

	hcf, err := httpcfg.NewHTTPCfg(http)
	if err != nil {
		return cfg, err
	}
	cfg.HTTP = *hcf
	cfg.BaseURL = *baseURLcfg.NewBaseURLCfg(baseURL)
	cfg.Log = *logscfg.NewLogsCfg()
	cfg.Storage = *storagecfg.NewStorageCfg(filePath)

	return cfg, nil
}

func parseHTTPServerFlags() (*httpcfg.HTTPCfg, *baseURLcfg.BaseURLCfg, *storagecfg.FileStorage) {
	addr := new(httpcfg.HTTPCfg)
	url := new(baseURLcfg.BaseURLCfg)
	filePath := new(storagecfg.FileStorage)

	flag.Var(addr, "a", "Net address host:port")
	flag.Var(url, "b", "base address of the resulting shortened URL http://localhost:3000")
	flag.Var(filePath, "f", "path to file storage")
	flag.Parse()

	return addr, url, filePath
}
