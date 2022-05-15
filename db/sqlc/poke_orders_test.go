package db

import (
	"context"
	"testing"

	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func mockOrderData(t *testing.T, user User, pokemon PokeProduct) PokeOrder {
	arg := InsertPokemonOrderDataParams{
		UserID:      user.ID,
		ProductID:   pokemon.ID,
		Quantity:    int32(util.RandomAmount()),
		TotalPrice:  util.RandomAmount(),
		OrderDetail: util.RandomRole(),
	}

	data, err := testQueries.InsertPokemonOrderData(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	require.Equal(t, arg.UserID, data.UserID)
	require.Equal(t, arg.ProductID, data.ProductID)
	require.Equal(t, arg.Quantity, data.Quantity)
	require.Equal(t, arg.TotalPrice, data.TotalPrice)
	require.Equal(t, arg.OrderDetail, data.OrderDetail)

	require.NotZero(t, data.ID)
	require.NotZero(t, data.CreatedAt)

	return data
}

func TestInsertPokemonOrderData(t *testing.T) {
	user := mockCreateUserAccount(t)
	poke := mockRandomData(t)
	mockOrderData(t, user, poke)
}

func TestListPokemonOrderData(t *testing.T) {
	user := mockCreateUserAccount(t)
	poke := mockRandomData(t)
	for i := 0; i < 10; i++ {
		mockOrderData(t, user, poke)
	}
	arg := ListPokemonOrderDataParams{
		UserID:    user.ID,
		ProductID: poke.ID,
		Limit:     5,
		Offset:    5,
	}

	orders, err := testQueries.ListPokemonOrderData(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, orders, 5)

	for _, order := range orders {
		require.NotEmpty(t, order)
	}
}

func TestCancelPokemonOrderData(t *testing.T) {
	user := mockCreateUserAccount(t)
	poke := mockRandomData(t)
	data := mockOrderData(t, user, poke)

	cancel, err := testQueries.CancelPokemonOrderData(context.Background(), data.ID)
	require.NoError(t, err)
	require.NotZero(t, cancel)

}
