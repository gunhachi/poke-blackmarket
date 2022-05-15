-- name: InsertPokemonOrderData :one
INSERT INTO poke_orders (
    user_id, product_id,quantity,total_price,order_detail
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListPokemonOrderData :many
SELECT * FROM poke_orders
-- WHERE 
--     user_id = $1 OR
--     product_id = $2
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: CancelPokemonOrderData :exec
DELETE FROM poke_orders
WHERE id = $1;

-- name: GetPokemonOrderData :one
SELECT * FROM poke_orders
WHERE id = $1 LIMIT 1;
