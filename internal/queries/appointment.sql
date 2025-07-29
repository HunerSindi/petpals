-- name: CreateAppointment :one
INSERT INTO appointments (
    user_id, clinic_id, pet_id, appointment_date, appointment_time, status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetAppointmentByID :one
SELECT * FROM appointments
WHERE id = $1 LIMIT 1;

-- name: ListAppointmentsByUserID :many
SELECT * FROM appointments
WHERE user_id = $1
ORDER BY appointment_date DESC, appointment_time DESC;

-- name: ListAppointmentsDetailsByUserID :many
SELECT
    a.id,
    a.appointment_date,
    a.appointment_time,
    a.status,
    a.created_at,
    u.id as user_id,
    u.first_name,
    u.last_name,
    u.phone,
    u.email,
    c.id as clinic_id,
    c.clinic_name,
    c.email as clinic_email,
    c.open_time,
    c.close_time,
    c.description as clinic_description,
    p.id as pet_id,
    p.name as pet_name,
    p.type as pet_type,
    p.birth_date as pet_birth_date
FROM
    appointments a
JOIN
    users u ON a.user_id = u.id
JOIN
    clinics c ON a.clinic_id = c.id
JOIN
    pets p ON a.pet_id = p.id
WHERE
    a.user_id = $1
ORDER BY
    a.appointment_date DESC, a.appointment_time DESC;

-- name: ListAppointmentsByClinicID :many
SELECT * FROM appointments
WHERE clinic_id = $1
ORDER BY appointment_date DESC, appointment_time DESC;

-- name: UpdateAppointment :one
UPDATE appointments
SET
    clinic_id = COALESCE(sqlc.narg(clinic_id), clinic_id),
    pet_id = COALESCE(sqlc.narg(pet_id), pet_id),
    appointment_date = COALESCE(sqlc.narg(appointment_date), appointment_date),
    appointment_time = COALESCE(sqlc.narg(appointment_time), appointment_time),
    status = COALESCE(sqlc.narg(status), status)
WHERE
    id = $1
RETURNING *;

-- name: DeleteAppointment :exec
DELETE FROM appointments
WHERE id = $1;