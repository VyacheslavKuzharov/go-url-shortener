package postgres

import (
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
	"time"
)

const (
	defaultAttempts = 5
	defaultTimeout  = time.Second
)

func RunMigrations(connectURL string, l *logger.Logger) {
	l.Info("Current storage is PostgreSQL. Trying to run migrations...")

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		m, err = migrate.New("file://migrations", connectURL)
		if err == nil {
			break
		}

		l.Info("Migrate: postgres is trying to connect, attempts left: %d", attempts)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		log.Fatalf("Migrate: postgres connect error: %s", err)
	}

	err = m.Up()
	defer m.Close()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Migrate: up error: %s", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		l.Info("Migrate: no change")
		return
	}

	l.Info("Migrate: up success")
}
