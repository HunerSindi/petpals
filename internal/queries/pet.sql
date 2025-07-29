-- name: CreatePet :one
INSERT INTO pets (
    uuid, user_id, name, type, birth_date
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPetByID :one
SELECT * FROM pets
WHERE id = $1 LIMIT 1;

-- name: ListPetsByUserID :many
SELECT * FROM pets
WHERE user_id = $1
ORDER BY id;

-- name: UpdatePet :one
UPDATE pets
SET
    name = COALESCE(sqlc.narg(name), name),
    type = COALESCE(sqlc.narg(type), type),
    birth_date = COALESCE(sqlc.narg(birth_date), birth_date)
WHERE
    id = $1
RETURNING *;

-- name: DeletePet :exec
DELETE FROM pets
WHERE id = $1;
