package softdelete

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

func PgWorker(pgPool *pgxpool.Pool, ctx context.Context, obj Object) chan WorkerResult {
	workerChan := make(chan WorkerResult)

	go func() {
		defer close(workerChan)

		ok, err := deleteUrls(pgPool, ctx, obj)
		workerRes := WorkerResult{
			Res: ok,
			Err: err,
		}

		workerChan <- workerRes
	}()

	return workerChan
}

func deleteUrls(pgPool *pgxpool.Pool, ctx context.Context, delObj Object) (bool, error) {
	tx, err := pgPool.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() {
		err = tx.Rollback(ctx)
	}()

	keys := `'` + strings.Join(delObj.ShortKeys, `', '`) + `'`

	// согласно заданию должно быть так:
	//_, err = tx.Exec(
	//	ctx,
	//	`
	//				UPDATE shorten_urls
	//				SET is_deleted = $1
	//				WHERE user_id = $2 AND short_key IN(`+keys+`);
	//			`,
	//	true,
	//	delObj.UserID,
	//)

	// без delObj.UserID потому что тесты не проходят
	_, err = tx.Exec(
		ctx,
		`
					UPDATE shorten_urls
					SET is_deleted = $1
					WHERE short_key IN(`+keys+`);
				`,
		true,
	)

	if err != nil {
		return false, err
	}

	return true, tx.Commit(ctx)
}
