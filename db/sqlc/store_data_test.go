package db

import (
	"context"
	"testing"

	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func mockOrderTx(t *testing.T, user User, pokemon PokeProduct) OrderTxResult {
	order := NewStore(testDB)

	result, err := order.OrderTx(context.Background(), OrderTxParams{
		UserID:    user.ID,
		ProductID: pokemon.ID,
		Quantity:  int32(util.RandomAmount()),
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user.ID, result.Order.UserID)
	require.Equal(t, pokemon.ID, result.Order.ProductID)
	return result
}

func TestOrdertx(t *testing.T) {
	user := mockCreateUserAccount(t)
	pokemon := mockRandomData(t)
	mockOrderTx(t, user, pokemon)
}

func TestCancelOrdertx(t *testing.T) {
	order := NewStore(testDB)

	user := mockCreateUserAccount(t)
	pokemon := mockRandomData(t)
	data := mockOrderTx(t, user, pokemon)

	result, err := order.CancelOrderTx(context.Background(), CancelOrderParam{
		ID: data.Order.ID,
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

}
