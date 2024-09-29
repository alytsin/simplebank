package db

import (
	"context"
)

// TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *TxStore) Transfer(ctx context.Context, arg *TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, &CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, &CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, &CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		// порядок выполнения важен. если есть две транзакции с такой последовательностью, то получаем deadlock
		// транзакция 1:
		// a. update ... where id = 1
		// b. update ... where id = 2
		// транзакция 2:
		// a. update ... where id = 2
		// b. update ... where id = 1
		// делаем обновление по порядку следования id и deadlock'ов не будет.
		if arg.FromAccountID > arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = store.transferMoney(ctx, q,
				arg.FromAccountID,
				arg.ToAccountID,
				arg.Amount,
			)
		} else {
			result.ToAccount, result.FromAccount, err = store.transferMoney(ctx, q,
				arg.ToAccountID,
				arg.FromAccountID,
				-arg.Amount,
			)
		}

		return nil
	})

	return result, err
}

func (store *TxStore) transferMoney(
	ctx context.Context,
	q *Queries,
	fromAccountID int64,
	toAccountID int64,
	amount int64,
) (src Account, dst Account, err error) {

	src, err = q.AddAccountBalance(ctx, &AddAccountBalanceParams{
		Amount: -amount,
		ID:     fromAccountID,
	})
	if err != nil {
		return
	}

	dst, err = q.AddAccountBalance(ctx, &AddAccountBalanceParams{
		Amount: amount,
		ID:     toAccountID,
	})
	if err != nil {
		return
	}

	return
}
