package postgres

import (
	"context"
	"errors"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/config"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/entity"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/httpapi"
	"github.com/VyacheslavKuzharov/go-url-shortener/internal/lib/random"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	uuid "github.com/satori/go.uuid"
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

func (pg *Pg) SaveURL(ctx context.Context, originalURL string) (string, error) {
	if originalURL == "" {
		return "", errors.New("originalURL can't be blank")
	}
	shortKey := random.GenShortKey()
	userID := ctx.Value(entity.CurrentUserID)

	_, err := pg.Pool.Exec(
		ctx,
		"INSERT INTO shorten_urls(short_key, original_url, user_id) VALUES ($1, $2, $3)",
		shortKey,
		originalURL,
		userID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			key, _ := pg.getShortKey(ctx, originalURL)
			return "", NewUniqueFieldErr(originalURL, key, err)
		}

		return "", err
	}

	return shortKey, nil
}

func (pg *Pg) SaveBatchURLs(ctx context.Context, urls []entity.ShortenURL) error {
	if len(urls) == 0 {
		return nil
	}

	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		err = tx.Rollback(ctx)
	}()

	var batch = &pgx.Batch{}
	userID := ctx.Value(entity.CurrentUserID)

	for _, u := range urls {
		batch.Queue(
			`INSERT INTO shorten_urls (short_key, original_url, user_id) VALUES ($1, $2, $3)`,
			u.ShortKey,
			u.OriginalURL,
			userID,
		)
	}

	res := tx.SendBatch(ctx, batch)
	err = res.Close()
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (pg *Pg) GetURL(ctx context.Context, key string) (string, error) {
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

func (pg *Pg) GetUserUrls(ctx context.Context, currentUserID uuid.UUID, cfg *config.Config) ([]*entity.CompletedURL, error) {
	var userURLs []*entity.CompletedURL

	rows, err := pg.Pool.Query(
		ctx,
		"SELECT short_key, original_url FROM shorten_urls WHERE user_id = $1",
		currentUserID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shortenURLs []entity.ShortenURL
	for rows.Next() {
		var su entity.ShortenURL

		err = rows.Scan(&su.ShortKey, &su.OriginalURL)
		if err != nil {
			return nil, err
		}

		shortenURLs = append(shortenURLs, su)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for _, v := range shortenURLs {
		urlItem := &entity.CompletedURL{
			ShortURL:    httpapi.FullShortenedURL(v.ShortKey, cfg),
			OriginalURL: v.OriginalURL,
		}
		userURLs = append(userURLs, urlItem)
	}

	return userURLs, nil
}

func (pg *Pg) Close() error {
	pg.Pool.Close()
	return nil
}

func (pg *Pg) Ping(ctx context.Context) error {
	if err := pg.Pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func (pg *Pg) getShortKey(ctx context.Context, originalURL string) (string, error) {
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
