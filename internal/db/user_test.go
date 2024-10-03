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
	require.Zero(t, user.PasswordChangedAt)

	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
