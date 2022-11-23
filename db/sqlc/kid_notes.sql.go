// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: kid_notes.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const createKidNote = `-- name: CreateKidNote :one
INSERT INTO kid_notes(note, kid_id, presence, has_meal, date)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, note, kid_id, presence, has_meal, date
`

type CreateKidNoteParams struct {
	Note     string     `json:"note"`
	KidID    int64      `json:"kidID"`
	Presence []Presence `json:"presence"`
	HasMeal  bool       `json:"hasMeal"`
	Date     time.Time  `json:"date"`
}

func (q *Queries) CreateKidNote(ctx context.Context, arg CreateKidNoteParams) (KidNote, error) {
	row := q.db.QueryRowContext(ctx, createKidNote,
		arg.Note,
		arg.KidID,
		pq.Array(arg.Presence),
		arg.HasMeal,
		arg.Date,
	)
	var i KidNote
	err := row.Scan(
		&i.ID,
		&i.Note,
		&i.KidID,
		pq.Array(&i.Presence),
		&i.HasMeal,
		&i.Date,
	)
	return i, err
}

const getAllKidNotes = `-- name: GetAllKidNotes :many
SELECT id, note, kid_id, presence, has_meal, date
FROM kid_notes
ORDER BY id
`

func (q *Queries) GetAllKidNotes(ctx context.Context) ([]KidNote, error) {
	rows, err := q.db.QueryContext(ctx, getAllKidNotes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []KidNote{}
	for rows.Next() {
		var i KidNote
		if err := rows.Scan(
			&i.ID,
			&i.Note,
			&i.KidID,
			pq.Array(&i.Presence),
			&i.HasMeal,
			&i.Date,
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

const getKidNote = `-- name: GetKidNote :one
SELECT id, note, kid_id, presence, has_meal, date
FROM kid_notes
WHERE id = $1
`

func (q *Queries) GetKidNote(ctx context.Context, id int64) (KidNote, error) {
	row := q.db.QueryRowContext(ctx, getKidNote, id)
	var i KidNote
	err := row.Scan(
		&i.ID,
		&i.Note,
		&i.KidID,
		pq.Array(&i.Presence),
		&i.HasMeal,
		&i.Date,
	)
	return i, err
}

const getKidNotesByPeriod = `-- name: GetKidNotesByPeriod :many
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
  AND date <= $2
`

type GetKidNotesByPeriodParams struct {
	Date   time.Time `json:"date"`
	Date_2 time.Time `json:"date2"`
}

type GetKidNotesByPeriodRow struct {
	ID         int64          `json:"id"`
	Note       string         `json:"note"`
	Date       time.Time      `json:"date"`
	Presence   []Presence     `json:"presence"`
	HasMeal    bool           `json:"hasMeal"`
	KidName    sql.NullString `json:"kidName"`
	KidSurname sql.NullString `json:"kidSurname"`
	KidID      sql.NullInt64  `json:"kidID"`
}

func (q *Queries) GetKidNotesByPeriod(ctx context.Context, arg GetKidNotesByPeriodParams) ([]GetKidNotesByPeriodRow, error) {
	rows, err := q.db.QueryContext(ctx, getKidNotesByPeriod, arg.Date, arg.Date_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetKidNotesByPeriodRow{}
	for rows.Next() {
		var i GetKidNotesByPeriodRow
		if err := rows.Scan(
			&i.ID,
			&i.Note,
			&i.Date,
			pq.Array(&i.Presence),
			&i.HasMeal,
			&i.KidName,
			&i.KidSurname,
			&i.KidID,
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

const updateKidNote = `-- name: UpdateKidNote :one
UPDATE kid_notes
SET note=$2,
    presence=$3,
    has_meal=$4
WHERE id = $1
RETURNING id, note, kid_id, presence, has_meal, date
`

type UpdateKidNoteParams struct {
	ID       int64      `json:"id"`
	Note     string     `json:"note"`
	Presence []Presence `json:"presence"`
	HasMeal  bool       `json:"hasMeal"`
}

func (q *Queries) UpdateKidNote(ctx context.Context, arg UpdateKidNoteParams) (KidNote, error) {
	row := q.db.QueryRowContext(ctx, updateKidNote,
		arg.ID,
		arg.Note,
		pq.Array(arg.Presence),
		arg.HasMeal,
	)
	var i KidNote
	err := row.Scan(
		&i.ID,
		&i.Note,
		&i.KidID,
		pq.Array(&i.Presence),
		&i.HasMeal,
		&i.Date,
	)
	return i, err
}
