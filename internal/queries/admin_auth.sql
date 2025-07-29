-- name: GetAdminUserByUsername :one
SELECT * FROM admin_users
WHERE username = $1 LIMIT 1;

-- name: CreateAdminUser :one
INSERT INTO admin_users (
    username, password
) VALUES (
    $1, $2
) RETURNING *;
