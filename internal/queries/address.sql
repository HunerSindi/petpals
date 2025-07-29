-- name: CreateUserAddress :one
INSERT INTO user_addresses (
    user_id, address_line1, address_line2, city, state, postal_code, is_default
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetUserAddressByID :one
SELECT * FROM user_addresses
WHERE id = $1 LIMIT 1;

-- name: ListUserAddressesByUserID :many
SELECT * FROM user_addresses
WHERE user_id = $1
ORDER BY id;

-- name: UpdateUserAddress :one
UPDATE user_addresses
SET
    address_line1 = COALESCE(sqlc.narg(address_line1), address_line1),
    address_line2 = COALESCE(sqlc.narg(address_line2), address_line2),
    city = COALESCE(sqlc.narg(city), city),
    state = COALESCE(sqlc.narg(state), state),
    postal_code = COALESCE(sqlc.narg(postal_code), postal_code),
    is_default = COALESCE(sqlc.narg(is_default), is_default)
WHERE
    id = $1
RETURNING *;

-- name: DeleteUserAddress :exec
DELETE FROM user_addresses
WHERE id = $1;
