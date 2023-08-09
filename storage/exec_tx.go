package storage

import (
	"context"
	"database/sql"
	"fmt"
)

func (store *SQLStorage) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := NewQueries(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
