-- name: CreateTeacher :one
INSERT INTO teachers(name, surname)
VALUES ($1, $2)
RETURNING *;

-- name: GetAllTeachers :many
SELECT * FROM teachers
ORDER BY name;

-- name: GetTeacher :one
SELECT * FROM teachers
WHERE id=$1
LIMIT 1;