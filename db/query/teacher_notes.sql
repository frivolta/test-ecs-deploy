-- name: CreateTeacherNote :one
INSERT INTO teacher_notes(note, teacher_id, date)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateNote :one
UPDATE teacher_notes
SET note=$2,
    date=$3
WHERE id = $1
RETURNING *;

-- name: GetTeacherNote :one
SELECT *
FROM teacher_notes
WHERE id = $1;

-- name: GetAllTeacherNotes :many
SELECT *
FROM teacher_notes
ORDER BY id;

-- name: GetTeacherNotesByDate :many
SELECT teacher_notes.id,
       note,
       date,
       teachers.name    AS "teacher_name",
       teachers.surname AS "teacher_surname",
       teachers.id      AS "teacher_id"
FROM teacher_notes
         LEFT JOIN teachers ON teacher_notes.teacher_id = teachers.id
WHERE date = $1;

-- name: GetTeacherNotesByPeriod :many
SELECT teacher_notes.id,
       note,
       date,
       teachers.name    AS "teacher_name",
       teachers.surname AS "teacher_surname",
       teachers.id      AS "teacher_id"
FROM teacher_notes
         LEFT JOIN teachers ON teacher_notes.teacher_id = teachers.id
WHERE date >= $1
  AND date <= $2
ORDER BY date DESC;




