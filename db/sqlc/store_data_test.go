package db

import (
	"context"
	"testing"

	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func TestOrdertx(t *testing.T) {
	order := NewStore(testDB)

	user := mockCreateUserAccount(t)
	pokemon := mockRandomData(t)

	result, err := order.OrderTx(context.Background(), OrderTxParams{
		UserID:    user.ID,
		ProductID: pokemon.ID,
		Quantity:  int32(util.RandomAmount()),
	})

	require.NoError(t, err)
	require.NotEmpty(t, result)

	require.Equal(t, user.ID, result.Order.UserID)
	require.Equal(t, pokemon.ID, result.Order.ProductID)

	// require.Error(t, err, errors.New("quantity exceed"))
}
