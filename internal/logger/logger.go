package logger

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"strings"
	"time"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(level string) *Logger {
	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(l)

	loggerOutput := zerolog.ConsoleWriter{Out: os.Stderr}
	logger := zerolog.New(loggerOutput).With().Timestamp().Logger()

	return &Logger{
		logger: &logger,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.log(msg, args...)
}

func (l *Logger) LogHTTPRequest(ww middleware.WrapResponseWriter, r *http.Request, start time.Time) {
	l.logger.Info().
		Int("status", ww.Status()).
		Str("method", r.Method).
		Str("path", r.URL.Path).
		Str("uri", r.RequestURI).
		Str("query", r.URL.RawQuery).
		Dur("duration", time.Since(start)).
		Int("bytes", ww.BytesWritten()).
		Msg("request completed")
}

func (l *Logger) log(msg string, args ...any) {
	if len(args) == 0 {
		l.logger.Info().Msg(msg)
	} else {
		l.logger.Info().Msgf(msg, args...)
	}
}
