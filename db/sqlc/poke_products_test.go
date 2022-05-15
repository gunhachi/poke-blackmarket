package db

import (
	"context"
	"testing"
	"time"

	"github.com/gunhachi/poke-blackmarket/util"
	"github.com/stretchr/testify/require"
)

func mockRandomData(t *testing.T) PokeProduct {
	arg := CreatePokemonDataParams{
		PokeName:  util.RandomUser(),
		Status:    util.RandomUser(),
		PokeStock: util.RandomAmount(),
		PokePrice: util.RandomAmount(),
	}

	data, err := testQueries.CreatePokemonData(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, data)

	require.Equal(t, arg.PokeName, data.PokeName)
	require.Equal(t, arg.Status, data.Status)
	require.Equal(t, arg.PokeStock, data.PokeStock)
	require.Equal(t, arg.PokePrice, data.PokePrice)

	require.NotZero(t, data.ID)
	require.NotZero(t, data.CreatedAt)

	return data
}

func TestCreatePokemonData(t *testing.T) {
	mockRandomData(t)
}

func TestGetPokemonData(t *testing.T) {
	data1 := mockRandomData(t)
	data2, err := testQueries.GetPokemonData(context.Background(), data1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, data2)

	require.Equal(t, data1.ID, data2.ID)
	require.Equal(t, data1.PokeName, data2.PokeName)
	require.Equal(t, data1.Status, data2.Status)
	require.Equal(t, data1.PokeStock, data2.PokeStock)
	require.Equal(t, data1.PokePrice, data2.PokePrice)
	require.WithinDuration(t, data1.CreatedAt, data2.CreatedAt, time.Second)
}

func TestAddPokemonStockData(t *testing.T) {
	data1 := mockRandomData(t)

	arg := AddPokemonStockDataParams{
		Amount: util.RandomAmount(),
		ID:     data1.ID,
	}

	data2, err := testQueries.AddPokemonStockData(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, data2)

	require.Equal(t, data1.ID, data2.ID)
	require.Equal(t, data1.PokeName, data2.PokeName)
	require.Equal(t, data1.Status, data2.Status)
	require.Equal(t, data2.PokeStock, data1.PokeStock+arg.Amount)
	require.Equal(t, data1.PokePrice, data2.PokePrice)
	require.WithinDuration(t, data1.CreatedAt, data2.CreatedAt, time.Second)
}

func TestDeductPokemonStockData(t *testing.T) {
	data1 := mockRandomData(t)

	arg := DeductPokemonStockDataParams{
		Amount: util.RandomAmount(),
		ID:     data1.ID,
	}

	data2, err := testQueries.DeductPokemonStockData(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, data2)

	require.Equal(t, data1.ID, data2.ID)
	require.Equal(t, data1.PokeName, data2.PokeName)
	require.Equal(t, data1.Status, data2.Status)
	require.Equal(t, data2.PokeStock, data1.PokeStock-arg.Amount)
	require.Equal(t, data1.PokePrice, data2.PokePrice)
	require.WithinDuration(t, data1.CreatedAt, data2.CreatedAt, time.Second)
}

func TestListPokemonData(t *testing.T) {
	for i := 0; i < 10; i++ {
		mockRandomData(t)
	}

	arg := ListPokemonDataParams{
		Limit:  5,
		Offset: 5,
	}

	data, err := testQueries.ListPokemonData(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, data, 5)

	for _, poke := range data {
		require.NotEmpty(t, poke)
	}
}

func TestUpdatePokemonStockData(t *testing.T) {
	data1 := mockRandomData(t)

	arg := UpdatePokemonDataParams{
		ID:        data1.ID,
		Status:    data1.Status,
		PokePrice: data1.PokePrice,
		PokeStock: data1.PokeStock,
	}

	data2, err := testQueries.UpdatePokemonData(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, data2)

	require.Equal(t, data1.ID, data2.ID)
	require.Equal(t, data1.PokeName, data2.PokeName)
	require.Equal(t, data1.Status, data2.Status)
	require.Equal(t, data1.PokePrice, data2.PokePrice)
	require.WithinDuration(t, data1.CreatedAt, data2.CreatedAt, time.Second)
}
