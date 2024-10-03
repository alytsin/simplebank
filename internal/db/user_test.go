package db

import (
	"context"
	"github.com/alytsin/simplebank/internal/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomUser(t *testing.T) *User {
	params := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: util.RandomOwner(),
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), params)

	require.Nil(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, params.Username, user.Username)
	require.Equal(t, params.HashedPassword, user.HashedPassword)
	require.Equal(t, params.FullName, user.FullName)
	require.Equal(t, params.Email, user.Email)
	require.True(t, user.PasswordChangedAt.IsZero())

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestDuplicateUser(t *testing.T) {
	user := createRandomUser(t)

	params := CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}

	dup, err := testQueries.CreateUser(context.Background(), params)
	require.Error(t, err)
	require.Empty(t, dup)

}

func TestGetUser(t *testing.T) {

	user := createRandomUser(t)

	found, err := testQueries.GetUser(context.Background(), user.Username)
	require.Nil(t, err)
	require.Equal(t, user.Username, found.Username)
	require.Equal(t, user.HashedPassword, found.HashedPassword)
	require.Equal(t, user.FullName, found.FullName)
	require.Equal(t, user.Email, found.Email)
	require.Equal(t, user.CreatedAt, found.CreatedAt)
	require.Equal(t, user.PasswordChangedAt, found.PasswordChangedAt)

}
