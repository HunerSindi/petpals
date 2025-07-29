-- name: CreateProduct :one
INSERT INTO products (
    category_id, name, description, price
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
    category_id = COALESCE(sqlc.narg(category_id), category_id),
    name = COALESCE(sqlc.narg(name), name),
    description = COALESCE(sqlc.narg(description), description),
    price = COALESCE(sqlc.narg(price), price)
WHERE
    id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: CreateProductImage :one
INSERT INTO product_images (
    product_id, img_url, is_primary
) VALUES (
    $1, $2, $3
) RETURNING *;
