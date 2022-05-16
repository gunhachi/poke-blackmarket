package db

import (
	"context"
	"testing"
	"time"

	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func mockCreateUserAccount(t *testing.T) User {
	account := mockCreateAccountLog(t)
	arg := CreateUserAccountParams{
		UserName: account.Username,
		UserRole: util.RandomRole(),
	}

	user, err := testQueries.CreateUserAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.UserRole, user.UserRole)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCrateUserAccount(t *testing.T) {
	mockCreateUserAccount(t)
}

func TestGetUserAccount(t *testing.T) {
	user1 := mockCreateUserAccount(t)
	user2, err := testQueries.GetUserAccount(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.UserRole, user2.UserRole)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)

}

func TestUpadteUserAccountRole(t *testing.T) {
	user1 := mockCreateUserAccount(t)

	arg := UpdateUserAccountRoleParams{
		ID:       user1.ID,
		UserRole: util.RandomRole(),
	}

	user2, err := testQueries.UpdateUserAccountRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, arg.UserRole, user2.UserRole)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestDeleteUserAccount(t *testing.T) {
	user1 := mockCreateUserAccount(t)
	err := testQueries.DeleteUserAccount(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := testQueries.GetUserAccount(context.Background(), user1.ID)
	require.Error(t, err)
	require.Empty(t, user2)
}

func TestListUserAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		mockCreateUserAccount(t)
	}

	arg := ListUserAccountParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := testQueries.ListUserAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)

	for _, user := range users {
		require.NotEmpty(t, user)
	}

}
