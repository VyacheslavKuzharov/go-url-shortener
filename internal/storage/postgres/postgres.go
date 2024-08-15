package postgres

import (
	"context"
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
	return "", nil
}

func (pg *Pg) GetURL(key string) (string, error) {
	return "", nil
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
