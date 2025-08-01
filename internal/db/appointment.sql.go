// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: appointment.sql

package db

import (
	"context"
	"database/sql"
)

const createAppointment = `-- name: CreateAppointment :one
INSERT INTO appointments (
    user_id, clinic_id, pet_id, appointment_date, appointment_time, status
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING id, user_id, clinic_id, pet_id, appointment_date, appointment_time, status, created_at
`

type CreateAppointmentParams struct {
	UserID          sql.NullInt64         `json:"user_id"`
	ClinicID        sql.NullInt64         `json:"clinic_id"`
	PetID           sql.NullInt64         `json:"pet_id"`
	AppointmentDate sql.NullTime          `json:"appointment_date"`
	AppointmentTime sql.NullTime          `json:"appointment_time"`
	Status          NullAppointmentStatus `json:"status"`
}

func (q *Queries) CreateAppointment(ctx context.Context, arg CreateAppointmentParams) (Appointment, error) {
	row := q.db.QueryRowContext(ctx, createAppointment,
		arg.UserID,
		arg.ClinicID,
		arg.PetID,
		arg.AppointmentDate,
		arg.AppointmentTime,
		arg.Status,
	)
	var i Appointment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ClinicID,
		&i.PetID,
		&i.AppointmentDate,
		&i.AppointmentTime,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAppointment = `-- name: DeleteAppointment :exec
DELETE FROM appointments
WHERE id = $1
`

func (q *Queries) DeleteAppointment(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAppointment, id)
	return err
}

const getAppointmentByID = `-- name: GetAppointmentByID :one
SELECT id, user_id, clinic_id, pet_id, appointment_date, appointment_time, status, created_at FROM appointments
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAppointmentByID(ctx context.Context, id int64) (Appointment, error) {
	row := q.db.QueryRowContext(ctx, getAppointmentByID, id)
	var i Appointment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ClinicID,
		&i.PetID,
		&i.AppointmentDate,
		&i.AppointmentTime,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}

const listAppointmentsByClinicID = `-- name: ListAppointmentsByClinicID :many
SELECT id, user_id, clinic_id, pet_id, appointment_date, appointment_time, status, created_at FROM appointments
WHERE clinic_id = $1
ORDER BY appointment_date DESC, appointment_time DESC
`

func (q *Queries) ListAppointmentsByClinicID(ctx context.Context, clinicID sql.NullInt64) ([]Appointment, error) {
	rows, err := q.db.QueryContext(ctx, listAppointmentsByClinicID, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Appointment{}
	for rows.Next() {
		var i Appointment
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ClinicID,
			&i.PetID,
			&i.AppointmentDate,
			&i.AppointmentTime,
			&i.Status,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAppointmentsByUserID = `-- name: ListAppointmentsByUserID :many
SELECT id, user_id, clinic_id, pet_id, appointment_date, appointment_time, status, created_at FROM appointments
WHERE user_id = $1
ORDER BY appointment_date DESC, appointment_time DESC
`

func (q *Queries) ListAppointmentsByUserID(ctx context.Context, userID sql.NullInt64) ([]Appointment, error) {
	rows, err := q.db.QueryContext(ctx, listAppointmentsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Appointment{}
	for rows.Next() {
		var i Appointment
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.ClinicID,
			&i.PetID,
			&i.AppointmentDate,
			&i.AppointmentTime,
			&i.Status,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAppointmentsDetailsByUserID = `-- name: ListAppointmentsDetailsByUserID :many
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
    a.appointment_date DESC, a.appointment_time DESC
`

type ListAppointmentsDetailsByUserIDRow struct {
	ID                int64                 `json:"id"`
	AppointmentDate   sql.NullTime          `json:"appointment_date"`
	AppointmentTime   sql.NullTime          `json:"appointment_time"`
	Status            NullAppointmentStatus `json:"status"`
	CreatedAt         sql.NullTime          `json:"created_at"`
	UserID            int64                 `json:"user_id"`
	FirstName         sql.NullString        `json:"first_name"`
	LastName          sql.NullString        `json:"last_name"`
	Phone             sql.NullString        `json:"phone"`
	Email             sql.NullString        `json:"email"`
	ClinicID          int64                 `json:"clinic_id"`
	ClinicName        sql.NullString        `json:"clinic_name"`
	ClinicEmail       sql.NullString        `json:"clinic_email"`
	OpenTime          sql.NullTime          `json:"open_time"`
	CloseTime         sql.NullTime          `json:"close_time"`
	ClinicDescription sql.NullString        `json:"clinic_description"`
	PetID             int64                 `json:"pet_id"`
	PetName           sql.NullString        `json:"pet_name"`
	PetType           sql.NullString        `json:"pet_type"`
	PetBirthDate      sql.NullTime          `json:"pet_birth_date"`
}

func (q *Queries) ListAppointmentsDetailsByUserID(ctx context.Context, userID sql.NullInt64) ([]ListAppointmentsDetailsByUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, listAppointmentsDetailsByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ListAppointmentsDetailsByUserIDRow{}
	for rows.Next() {
		var i ListAppointmentsDetailsByUserIDRow
		if err := rows.Scan(
			&i.ID,
			&i.AppointmentDate,
			&i.AppointmentTime,
			&i.Status,
			&i.CreatedAt,
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.Phone,
			&i.Email,
			&i.ClinicID,
			&i.ClinicName,
			&i.ClinicEmail,
			&i.OpenTime,
			&i.CloseTime,
			&i.ClinicDescription,
			&i.PetID,
			&i.PetName,
			&i.PetType,
			&i.PetBirthDate,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAppointment = `-- name: UpdateAppointment :one
UPDATE appointments
SET
    clinic_id = COALESCE($2, clinic_id),
    pet_id = COALESCE($3, pet_id),
    appointment_date = COALESCE($4, appointment_date),
    appointment_time = COALESCE($5, appointment_time),
    status = COALESCE($6, status)
WHERE
    id = $1
RETURNING id, user_id, clinic_id, pet_id, appointment_date, appointment_time, status, created_at
`

type UpdateAppointmentParams struct {
	ID              int64                 `json:"id"`
	ClinicID        sql.NullInt64         `json:"clinic_id"`
	PetID           sql.NullInt64         `json:"pet_id"`
	AppointmentDate sql.NullTime          `json:"appointment_date"`
	AppointmentTime sql.NullTime          `json:"appointment_time"`
	Status          NullAppointmentStatus `json:"status"`
}

func (q *Queries) UpdateAppointment(ctx context.Context, arg UpdateAppointmentParams) (Appointment, error) {
	row := q.db.QueryRowContext(ctx, updateAppointment,
		arg.ID,
		arg.ClinicID,
		arg.PetID,
		arg.AppointmentDate,
		arg.AppointmentTime,
		arg.Status,
	)
	var i Appointment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.ClinicID,
		&i.PetID,
		&i.AppointmentDate,
		&i.AppointmentTime,
		&i.Status,
		&i.CreatedAt,
	)
	return i, err
}
