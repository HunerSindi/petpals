-- name: CreateClinic :one
INSERT INTO clinics (
    first_name, last_name, clinic_name, email, password, open_time, close_time, description
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: DeleteClinic :exec
DELETE FROM clinics
WHERE id = $1;
