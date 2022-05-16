// Code generated by sqlc. DO NOT EDIT.
// source: accounts.sql

package db

import (
	"context"
)

const createAccountLog = `-- name: CreateAccountLog :one
INSERT INTO accounts (
    username, hashed_password, full_name
) VALUES (
    $1, $2, $3
) RETURNING username, hashed_password, full_name, created_at, password_changet_at
`

type CreateAccountLogParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
}

func (q *Queries) CreateAccountLog(ctx context.Context, arg CreateAccountLogParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccountLog, arg.Username, arg.HashedPassword, arg.FullName)
	var i Account
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.CreatedAt,
		&i.PasswordChangetAt,
	)
	return i, err
}

const getAccountLog = `-- name: GetAccountLog :one
SELECT username, hashed_password, full_name, created_at, password_changet_at FROM accounts
WHERE username = $1 LIMIT 1
`

func (q *Queries) GetAccountLog(ctx context.Context, username string) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountLog, username)
	var i Account
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.CreatedAt,
		&i.PasswordChangetAt,
	)
	return i, err
}