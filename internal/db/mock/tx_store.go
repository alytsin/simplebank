package mock

import (
	"context"
	"github.com/alytsin/simplebank/internal/db"
	"github.com/stretchr/testify/mock"
)

type TxStore struct {
	mock.Mock
}

func (s *TxStore) Transfer(ctx context.Context, arg *db.TransferTxParams) (*db.TransferTxResult, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).(*db.TransferTxResult), args.Error(1)
}

func (s *TxStore) AddAccountBalance(ctx context.Context, arg *db.AddAccountBalanceParams) (*db.Account, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).(*db.Account), args.Error(1)
}

func (s *TxStore) CreateAccount(ctx context.Context, arg *db.CreateAccountParams) (*db.Account, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).(*db.Account), args.Error(1)
}

func (s *TxStore) CreateEntry(ctx context.Context, arg *db.CreateEntryParams) (*db.Entry, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).(*db.Entry), args.Error(1)
}

func (s *TxStore) CreateTransfer(ctx context.Context, arg *db.CreateTransferParams) (*db.Transfer, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).(*db.Transfer), args.Error(1)
}

func (s *TxStore) DeleteAccount(ctx context.Context, id int64) error {
	args := s.Called(ctx, id)
	return args.Error(0)
}

func (s *TxStore) GetAccount(ctx context.Context, id int64) (*db.Account, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*db.Account), args.Error(1)
}

func (s *TxStore) GetAccountForUpdate(ctx context.Context, id int64) (*db.Account, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*db.Account), args.Error(1)
}

func (s *TxStore) GetEntry(ctx context.Context, id int64) (*db.Entry, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*db.Entry), args.Error(1)
}

func (s *TxStore) GetTransfer(ctx context.Context, id int64) (*db.Transfer, error) {
	args := s.Called(ctx, id)
	return args.Get(0).(*db.Transfer), args.Error(1)
}

func (s *TxStore) ListAccounts(ctx context.Context, arg *db.ListAccountsParams) ([]*db.Account, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).([]*db.Account), args.Error(1)
}

func (s *TxStore) ListEntries(ctx context.Context, arg *db.ListEntriesParams) ([]*db.Entry, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).([]*db.Entry), args.Error(1)
}

func (s *TxStore) ListTransfers(ctx context.Context, arg *db.ListTransfersParams) ([]*db.Transfer, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).([]*db.Transfer), args.Error(1)
}

func (s *TxStore) UpdateAccount(ctx context.Context, arg *db.UpdateAccountParams) (*db.Account, error) {
	args := s.Called(ctx, arg)
	return args.Get(0).(*db.Account), args.Error(1)
}
