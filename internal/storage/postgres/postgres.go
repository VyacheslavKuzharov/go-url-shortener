package postgres

import (
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/random"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pg struct {
	Pool *pgxpool.Pool
}

func New(connectURL string) (*Pg, error) {
	connPool, err := pgxpool.New(context.Background(), connectURL)
	if err != nil {
		return nil, err
	}

	return &Pg{
		Pool: connPool,
	}, nil
}

func (pg *Pg) SaveURL(originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("originalURL can't be blank")
	}
	ctx := context.Background()
	shortKey := random.GenShortKey()

	_, err := pg.Pool.Exec(
		ctx,
		"INSERT INTO shorten_urls(short_key, original_url) VALUES ($1, $2)",
		shortKey,
		originalURL,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			key, _ := pg.getShortKey(originalURL)
			return "", NewUniqueFieldErr(originalURL, key, err)
		}

		return "", err
	}

	return shortKey, nil
}

func (pg *Pg) SaveBatchURLs(ctx context.Context, urls *[]entity.ShortenURL) error {
	resolvedURLs := *urls

	if len(resolvedURLs) == 0 {
		return nil
	}

	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var batch = &pgx.Batch{}

	for _, u := range resolvedURLs {
		batch.Queue(
			`INSERT INTO shorten_urls (short_key, original_url) VALUES ($1, $2)`,
			u.ShortKey,
			u.OriginalURL,
		)
	}

	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (pg *Pg) GetURL(key string) (string, error) {
	ctx := context.Background()
	var originalURL string

	row := pg.Pool.QueryRow(
		ctx,
		"SELECT original_url FROM shorten_urls WHERE short_key = $1",
		key,
	)
	err := row.Scan(&originalURL)
	if err != nil {
		return "", errors.New("shortKey not found")
	}

	return originalURL, nil
}

func (pg *Pg) Close() error {
	pg.Pool.Close()
	return nil
}

func (pg *Pg) Ping() error {
	ctx := context.Background()

	conn, err := pg.Pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	err = conn.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (pg *Pg) getShortKey(originalURL string) (string, error) {
	ctx := context.Background()
	var shortKey string

	row := pg.Pool.QueryRow(
		ctx,
		"SELECT short_key FROM shorten_urls WHERE original_url = $1",
		originalURL,
	)
	err := row.Scan(&shortKey)
	if err != nil {
		return "", errors.New("shortKey not found")
	}

	return shortKey, nil
}
