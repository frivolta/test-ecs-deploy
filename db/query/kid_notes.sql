-- name: CreateKidNote :one
INSERT INTO kid_notes(note, kid_id, presence, has_meal, date)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateKidNote :one
UPDATE kid_notes
SET note=$2,
    presence=$3,
    has_meal=$4
WHERE id = $1
RETURNING *;

-- name: GetKidNote :one
SELECT *
FROM kid_notes
WHERE id = $1;

-- name: GetAllKidNotes :many
SELECT *
FROM kid_notes
ORDER BY id;

-- name: GetKidNotesByPeriod :many
SELECT kid_notes.id,
       kid_notes.note,
       kid_notes.date,
       kid_notes.presence,
       kid_notes.has_meal,
       kids.name    AS "kid_name",
       kids.surname AS "kid_surname",
       kids.id       AS "kid_id"
FROM kid_notes
         LEFT JOIN kids ON kid_notes.kid_id = kids.id
WHERE date >= $1
  AND date <= $2;