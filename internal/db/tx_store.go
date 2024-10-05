package db

import (
	"context"
	"database/sql"
	"fmt"
)

type TxStoreInterface interface {
	Querier
	Transfer(ctx context.Context, arg TransferTxParams) (*TransferTxResult, error)
}

type TxStore struct {
	*Queries
	db *sql.DB
}

func NewTxStore(db *sql.DB) TxStoreInterface {
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
			return fmt.Errorf("transaction err: %v, rollback err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
