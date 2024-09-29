package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TxStore struct {
	*Queries
	db *sql.DB
}

func NewTxStore(db *sql.DB) *TxStore {
	return &TxStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *TxStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err = fn(New(tx)); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
