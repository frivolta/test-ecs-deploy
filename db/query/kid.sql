-- name: CreateKid :one
INSERT INTO kids(name, surname)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllKids :many
SELECT * FROM kids
ORDER BY name;

-- name: GetKid :one
SELECT * FROM kids
WHERE id=$1
LIMIT 1;