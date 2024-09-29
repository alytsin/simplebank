package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type transferTxPair struct {
	R TransferTxResult
	E error
}

func TestTransferTx(t *testing.T) {

	c := make(chan transferTxPair)
	txStore := NewTxStore(testDb)

	from := createRandomAccount(t)
	to := createRandomAccount(t)

	//fmt.Println(">> before:", from.Balance, to.Balance)

	times := 10
	amount := int64(10)

	for i := 0; i < times; i++ {
		go func() {
			result, err := txStore.Transfer(context.Background(), &TransferTxParams{
				FromAccountID: from.ID,
				ToAccountID:   to.ID,
				Amount:        amount,
			})
			c <- transferTxPair{result, err}
		}()
	}

	for i := 0; i < times; i++ {
		result := <-c

		require.NotEmpty(t, result.R)
		require.Nil(t, result.E)

		tr := result.R.Transfer
		require.NotEmpty(t, tr)
		require.NotZero(t, tr.ID)
		require.NotZero(t, tr.CreatedAt)
		require.Equal(t, tr.Amount, amount)
		require.Equal(t, tr.FromAccountID, from.ID)
		require.Equal(t, tr.ToAccountID, to.ID)

		fe := result.R.FromEntry
		require.NotEmpty(t, fe)
		require.NotZero(t, fe.ID)
		require.NotZero(t, fe.CreatedAt)
		require.Equal(t, fe.AccountID, from.ID)
		require.Equal(t, fe.Amount, -amount)

		te := result.R.ToEntry
		require.NotEmpty(t, te)
		require.NotZero(t, te.ID)
		require.NotZero(t, te.CreatedAt)
		require.Equal(t, te.AccountID, to.ID)
		require.Equal(t, te.Amount, amount)

		fa := result.R.FromAccount
		require.NotEmpty(t, fa)
		require.NotZero(t, fa.ID)
		require.NotZero(t, fa.CreatedAt)

		ta := result.R.ToAccount
		require.NotEmpty(t, ta)
		require.NotZero(t, ta.ID)
		require.NotZero(t, ta.CreatedAt)

		//fmt.Println(">> tx:", fa.Balance, ta.Balance)

		require.Equal(t,
			from.Balance-fa.Balance,
			ta.Balance-to.Balance,
		)

	}

	ff, err := txStore.GetAccount(context.Background(), from.ID)
	require.Nil(t, err)
	require.NotEmpty(t, ff)

	tt, err := txStore.GetAccount(context.Background(), to.ID)
	require.Nil(t, err)
	require.NotEmpty(t, tt)

	//fmt.Println(">> after:", ff.Balance, tt.Balance)

	require.Equal(t, ff.Balance, from.Balance-int64(times)*amount)
	require.Equal(t, tt.Balance, to.Balance+int64(times)*amount)

}

func TestTransferTwoWaysTx(t *testing.T) {

	c := make(chan error)
	txStore := NewTxStore(testDb)

	from := createRandomAccount(t)
	to := createRandomAccount(t)

	times := 10
	amount := int64(10)

	for i := 0; i < times; i++ {
		go func() {

			fId := from.ID
			tId := to.ID

			if i%2 == 1 {
				fId = to.ID
				tId = from.ID
			}

			ctx := context.Background()
			_, err := txStore.Transfer(ctx, &TransferTxParams{
				FromAccountID: fId,
				ToAccountID:   tId,
				Amount:        amount,
			})
			c <- err
		}()
	}

	for i := 0; i < times; i++ {
		require.Nil(t, <-c)
	}

	ff, err := txStore.GetAccount(context.Background(), from.ID)
	require.Nil(t, err)
	require.NotEmpty(t, ff)

	tt, err := txStore.GetAccount(context.Background(), to.ID)
	require.Nil(t, err)
	require.NotEmpty(t, tt)

	require.Equal(t, ff.Balance, from.Balance)
	require.Equal(t, tt.Balance, to.Balance)

}
