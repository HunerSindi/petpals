-- name: CreateClinicImage :one
INSERT INTO clinic_images (
    clinic_id, img_url
) VALUES (
    $1, $2
) RETURNING *;

-- name: ListClinicImagesByClinicID :many
SELECT * FROM clinic_images
WHERE clinic_id = $1
ORDER BY id;
