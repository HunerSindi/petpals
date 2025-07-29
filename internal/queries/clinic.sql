-- name: ListClinics :many
SELECT * FROM clinics
ORDER BY id;

-- name: GetClinicByID :one
SELECT * FROM clinics
WHERE id = $1 LIMIT 1;

-- name: GetClinicAvailableSlots :many
SELECT * FROM appointments
WHERE clinic_id = $1 AND appointment_date >= NOW()
ORDER BY appointment_date, appointment_time;

-- name: UpdateClinicSchedule :one
UPDATE clinics
SET
    open_time = COALESCE(sqlc.narg(open_time), open_time),
    close_time = COALESCE(sqlc.narg(close_time), close_time)
WHERE
    id = $1
RETURNING *;