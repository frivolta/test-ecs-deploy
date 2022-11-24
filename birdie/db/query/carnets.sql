-- name: GetAllCarnets :many
SELECT carnets.id,
       date,
       quantity,
       kids.name    AS "kid_name",
       kids.surname AS "kid_surname",
       kids.id      AS "kid_id"
FROM carnets
         LEFT JOIN kids ON carnets.kid_id = kids.id
ORDER BY id;

-- name: CreateCarnet :one
INSERT INTO carnets(date, quantity, kid_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateCarnet :one
UPDATE carnets
SET date=$2,
    quantity=$3
WHERE id = $1
RETURNING *;

-- name: GetCarnet :one
SELECT * FROM carnets
WHERE id=$1
LIMIT 1;

