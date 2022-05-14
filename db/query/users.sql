-- name: CreateUserAccount :one
INSERT INTO users (
	user_name,user_role
) VALUES (
	$1, $2
) RETURNING *;

-- name: GetUserAccount :one
SELECT * FROM users 
WHERE id = $1 LIMIT 1;

-- name: ListUserAccount :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUserAccountRole :one
UPDATE users
SET user_role = $2
WHERE id = $1
RETURNING *;

-- name: DeleteUserAccount :exec
DELETE FROM users
WHERE id = $1;