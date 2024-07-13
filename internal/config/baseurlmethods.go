package config

func (b *BaseURLCfg) String() string {
	return b.Addr
}

func (b *BaseURLCfg) Set(s string) error {
	b.Addr = s
	return nil
}
