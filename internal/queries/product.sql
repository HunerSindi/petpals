-- name: ListProductsByCategory :many
SELECT * FROM products
WHERE category_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: GetProduct :one
SELECT * FROM products
WHERE id = $1 LIMIT 1;

-- name: ListProductImages :many
SELECT * FROM product_images
WHERE product_id = $1
ORDER BY is_primary DESC, id;
