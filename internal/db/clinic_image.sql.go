// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: clinic_image.sql

package db

import (
	"context"
	"database/sql"
)

const createClinicImage = `-- name: CreateClinicImage :one
INSERT INTO clinic_images (
    clinic_id, img_url
) VALUES (
    $1, $2
) RETURNING id, clinic_id, img_url
`

type CreateClinicImageParams struct {
	ClinicID sql.NullInt64 `json:"clinic_id"`
	ImgUrl   string        `json:"img_url"`
}

func (q *Queries) CreateClinicImage(ctx context.Context, arg CreateClinicImageParams) (ClinicImage, error) {
	row := q.db.QueryRowContext(ctx, createClinicImage, arg.ClinicID, arg.ImgUrl)
	var i ClinicImage
	err := row.Scan(&i.ID, &i.ClinicID, &i.ImgUrl)
	return i, err
}

const listClinicImagesByClinicID = `-- name: ListClinicImagesByClinicID :many
SELECT id, clinic_id, img_url FROM clinic_images
WHERE clinic_id = $1
ORDER BY id
`

func (q *Queries) ListClinicImagesByClinicID(ctx context.Context, clinicID sql.NullInt64) ([]ClinicImage, error) {
	rows, err := q.db.QueryContext(ctx, listClinicImagesByClinicID, clinicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ClinicImage{}
	for rows.Next() {
		var i ClinicImage
		if err := rows.Scan(&i.ID, &i.ClinicID, &i.ImgUrl); err != nil {
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
