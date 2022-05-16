package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provided functions to exec db query
type Store interface {
	Querier
	OrderTx(ctx context.Context, arg OrderTxParams) (OrderTxResult, error)
	CancelOrderTx(ctx context.Context, arg CancelOrderParam) (string, error)
}

// Store provided functions to exec db query
type SQLStore struct {
	*Queries
	db *sql.DB
}

// Create New Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

// OrderTxParams contains input parameter of the transaction
type OrderTxParams struct {
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type OrderTxResult struct {
	Order PokeOrder `json:"pokeorder"`
}

type CancelOrderParam struct {
	ID int64 `json:"id"`
}

// OrderTx perform Order transaction of pokemon and put it into table poke_orders
// It creates the order, add data in poke order, and update the pokemon stock based on pokemon id
func (store *SQLStore) OrderTx(ctx context.Context, arg OrderTxParams) (OrderTxResult, error) {
	var result OrderTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		getPokeData, err := q.GetPokemonData(ctx, arg.ProductID)
		if err != nil {
			return err
		}

		result.Order, err = q.InsertPokemonOrderData(ctx, InsertPokemonOrderDataParams{
			UserID:      arg.UserID,
			ProductID:   arg.ProductID,
			Quantity:    arg.Quantity,
			TotalPrice:  int64(arg.Quantity) * getPokeData.PokePrice,
			OrderDetail: "selling",
		})
		if err != nil {
			return err
		}

		_, err = q.DeductPokemonStockData(ctx, DeductPokemonStockDataParams{
			ID:     arg.ProductID,
			Amount: int64(arg.Quantity),
		})
		if err != nil {
			return err
		}

		return err

	})

	return result, err
}

// CancelOrderTx perform cancellation transaction of pokemon and return it stock data into table poke_orders
// It delete in poke order, and update the pokemon stock based on pokemon id
func (store *SQLStore) CancelOrderTx(ctx context.Context, arg CancelOrderParam) (string, error) {
	var result string

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		orderData, err := q.GetPokemonOrderData(ctx, arg.ID)
		if err != nil {
			return err
		}

		_, err = q.AddPokemonStockData(ctx, AddPokemonStockDataParams{
			ID:     orderData.ProductID,
			Amount: int64(orderData.Quantity),
		})
		if err != nil {
			result = "Error return stock"
			return err
		}

		return err

	})
	result = "stock returned"

	return result, err
}
