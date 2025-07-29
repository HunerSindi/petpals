-- name: CreateOrder :one
INSERT INTO orders (
    user_id, total_amount, status, delivery_address
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 LIMIT 1;

-- name: ListOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY order_date DESC;

-- name: UpdateOrderStatus :one
UPDATE orders
SET
    status = $2
WHERE
    id = $1
RETURNING *;
