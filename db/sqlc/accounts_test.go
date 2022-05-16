package db

import (
	"context"
	"testing"
	"time"

	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func mockCreateAccountLog(t *testing.T) Account {
	hashedPasswd, err := util.HashPassword(util.RandomString(5))
	require.NoError(t, err)

	arg := CreateAccountLogParams{
		Username:       util.RandomString(5),
		HashedPassword: hashedPasswd,
		FullName:       util.RandomString(10),
	}

	acc, err := testQueries.CreateAccountLog(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, acc)

	require.Equal(t, arg.Username, acc.Username)
	require.Equal(t, arg.HashedPassword, acc.HashedPassword)
	require.Equal(t, arg.FullName, acc.FullName)

	require.True(t, acc.PasswordChangetAt.Time.IsZero())
	require.NotZero(t, acc.CreatedAt)

	return acc
}

func TestCrateAccountLog(t *testing.T) {
	mockCreateAccountLog(t)
}

func TestGetUserAccountLog(t *testing.T) {
	user1 := mockCreateAccountLog(t)
	user2, err := testQueries.GetAccountLog(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.WithinDuration(t, user1.PasswordChangetAt.Time, user2.PasswordChangetAt.Time, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}
