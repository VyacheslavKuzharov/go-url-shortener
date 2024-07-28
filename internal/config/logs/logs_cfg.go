package config

type LogLevel string

const (
	InfoLevel  LogLevel = "info"
	ErrorLevel LogLevel = "error"
	WarnLevel  LogLevel = "warn"
)

type LogCfg struct {
	Level LogLevel
}

func NewLogsCfg() *LogCfg {
	return &LogCfg{
		Level: InfoLevel,
	}
}
