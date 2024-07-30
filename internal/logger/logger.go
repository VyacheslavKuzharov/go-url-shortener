package logger

import (
	logscfg "github.com/VyacheslavKuzharov/go-url-shortener/internal/config/logs"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	logger *zerolog.Logger
}

func New(level logscfg.LogLevel) *Logger {
	var l zerolog.Level

	switch level {
	case logscfg.ErrorLevel:
		l = zerolog.ErrorLevel
	case logscfg.WarnLevel:
		l = zerolog.WarnLevel
	case logscfg.InfoLevel:
		l = zerolog.InfoLevel
	default:
		l = zerolog.DebugLevel
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
