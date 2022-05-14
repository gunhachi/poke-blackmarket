// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0

package db

import (
	"time"
)

type PokeOrder struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"user_id"`
	ProductID int64 `json:"product_id"`
	// must be positive
	Quantity int32 `json:"quantity"`
	// must be positive
	TotalPrice  int64     `json:"total_price"`
	OrderDetail string    `json:"order_detail"`
	CreatedAt   time.Time `json:"created_at"`
}

type PokeProduct struct {
	ID       int64  `json:"id"`
	PokeName string `json:"poke_name"`
	Status   string `json:"status"`
	// must be positive
	PokePrice int64 `json:"poke_price"`
	// must be positive
	PokeStock int64     `json:"poke_stock"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	ID        int64     `json:"id"`
	UserName  string    `json:"user_name"`
	UserRole  string    `json:"user_role"`
	CreatedAt time.Time `json:"created_at"`
}
