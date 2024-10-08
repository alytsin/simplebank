package db

import (
	"context"
	"database/sql"
	"github.com/alytsin/simplebank/internal/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) *Account {

	user := createRandomUser(t)

	params := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), params)
	require.Nil(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, params.Owner, account.Owner)
	require.Equal(t, params.Balance, account.Balance)
	require.Equal(t, params.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
	//account := createRandomAccount(t)
	//_ = testQueries.DeleteAccount(context.Background(), account.ID)
}

func TestGetAccount(t *testing.T) {
	created := createRandomAccount(t)

	found, err := testQueries.GetAccount(context.Background(), created.ID)
	require.Nil(t, err)
	require.NotEmpty(t, found)

	require.Equal(t, created.ID, found.ID)
	require.Equal(t, created.Owner, found.Owner)
	require.Equal(t, created.Balance, found.Balance)
	require.Equal(t, created.Currency, found.Currency)
	require.Equal(t, created.CreatedAt, found.CreatedAt)

	//_ = testQueries.DeleteAccount(context.Background(), created.ID)
}

func TestUpdateAccount(t *testing.T) {
	created := createRandomAccount(t)

	update := UpdateAccountParams{
		ID:      created.ID,
		Balance: util.RandomMoney(),
	}

	updated, err := testQueries.UpdateAccount(context.Background(), update)
	require.Nil(t, err)
	require.NotEmpty(t, updated)

	require.Equal(t, created.ID, updated.ID)
	require.Equal(t, created.Owner, updated.Owner)
	require.Equal(t, update.Balance, updated.Balance)
	require.Equal(t, created.Currency, updated.Currency)
	require.Equal(t, created.CreatedAt, updated.CreatedAt)

	//_ = testQueries.DeleteAccount(context.Background(), created.ID)
}

func TestDuplicateAccount(t *testing.T) {

	created := createRandomAccount(t)

	dup, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    created.Owner,
		Currency: created.Currency,
		Balance:  0,
	})

	require.Error(t, err)
	require.Empty(t, dup)
}

func TestDeleteAccount(t *testing.T) {
	created := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), created.ID)
	require.Nil(t, err)

	found, err := testQueries.GetAccount(context.Background(), created.ID)
	require.Error(t, err)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.Empty(t, found)
}

func TestListAccounts(t *testing.T) {

	accounts := make([]*Account, 10)
	for i := 0; i < 10; i++ {
		accounts[i] = createRandomAccount(t)
	}

	for _, account := range accounts {
		found, err := testQueries.ListAccounts(context.Background(), ListAccountsParams{
			Owner:  account.Owner,
			Limit:  5,
			Offset: 0,
		})
		require.Nil(t, err)
		require.NotEmpty(t, found)
		require.Greater(t, len(found), 0)
	}

}
