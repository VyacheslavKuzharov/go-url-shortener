package config

import "os"

type BaseURLCfg struct {
	Addr string
}

func NewBaseURLCfg(flagBaseURL *BaseURLCfg) *BaseURLCfg {
	if flagBaseURL.Addr != "" {
		return flagBaseURL
	}

	return &BaseURLCfg{
		Addr: os.Getenv("BASE_URL"),
	}
}

func (b *BaseURLCfg) String() string {
	return b.Addr
}

func (b *BaseURLCfg) Set(s string) error {
	b.Addr = s
	return nil
}
