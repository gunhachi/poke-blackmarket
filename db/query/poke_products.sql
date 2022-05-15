-- name: CreatePokemonData :one
INSERT INTO poke_products (
    poke_name,status,poke_price,poke_stock
) VALUES (
	$1, $2, $3, $4
) RETURNING *;

-- name: GetPokemonData :one
SELECT * FROM poke_products
WHERE id = $1 LIMIT 1;

-- name: DeductPokemonStockData :one
UPDATE poke_products
SET poke_stock = poke_stock - sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: AddPokemonStockData :one
UPDATE poke_products
SET poke_stock = poke_stock + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: ListPokemonData :many
SELECT * FROM poke_products
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePokemonData :one
UPDATE poke_products
SET status = $2, poke_price = $3, poke_stock = $4
WHERE id = $1
RETURNING *;

