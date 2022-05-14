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