-- name: GetClinicByEmail :one
SELECT * FROM clinics
WHERE email = $1 LIMIT 1;

-- name: UpdateClinic :one
UPDATE clinics
SET
    first_name = COALESCE(sqlc.narg(first_name), first_name),
    last_name = COALESCE(sqlc.narg(last_name), last_name),
    clinic_name = COALESCE(sqlc.narg(clinic_name), clinic_name),
    email = COALESCE(sqlc.narg(email), email),
    password = COALESCE(sqlc.narg(password), password),
    open_time = COALESCE(sqlc.narg(open_time), open_time),
    close_time = COALESCE(sqlc.narg(close_time), close_time),
    description = COALESCE(sqlc.narg(description), description)
WHERE
    id = $1
RETURNING *;
