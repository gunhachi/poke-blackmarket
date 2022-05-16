-- name: CreateAccountLog :one
INSERT INTO accounts (
    username, hashed_password, full_name
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetAccountLog :one
SELECT * FROM accounts
WHERE username = $1 LIMIT 1;  