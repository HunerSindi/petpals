-- name: ListAllOrders :many
SELECT * FROM orders
ORDER BY order_date DESC;
