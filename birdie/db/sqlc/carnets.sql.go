// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: carnets.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const createCarnet = `-- name: CreateCarnet :one
INSERT INTO carnets(date, quantity, kid_id)
VALUES ($1, $2, $3)
RETURNING id, date, quantity, kid_id
`

type CreateCarnetParams struct {
	Date     time.Time `json:"date"`
	Quantity int32     `json:"quantity"`
	KidID    int64     `json:"kidID"`
}

func (q *Queries) CreateCarnet(ctx context.Context, arg CreateCarnetParams) (Carnet, error) {
	row := q.db.QueryRowContext(ctx, createCarnet, arg.Date, arg.Quantity, arg.KidID)
	var i Carnet
	err := row.Scan(
		&i.ID,
		&i.Date,
		&i.Quantity,
		&i.KidID,
	)
	return i, err
}

const getAllCarnets = `-- name: GetAllCarnets :many
SELECT carnets.id,
       date,
       quantity,
       kids.name    AS "kid_name",
       kids.surname AS "kid_surname",
       kids.id      AS "kid_id"
FROM carnets
         LEFT JOIN kids ON carnets.kid_id = kids.id
ORDER BY id
`

type GetAllCarnetsRow struct {
	ID         int64          `json:"id"`
	Date       time.Time      `json:"date"`
	Quantity   int32          `json:"quantity"`
	KidName    sql.NullString `json:"kidName"`
	KidSurname sql.NullString `json:"kidSurname"`
	KidID      sql.NullInt64  `json:"kidID"`
}

func (q *Queries) GetAllCarnets(ctx context.Context) ([]GetAllCarnetsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllCarnets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllCarnetsRow{}
	for rows.Next() {
		var i GetAllCarnetsRow
		if err := rows.Scan(
			&i.ID,
			&i.Date,
			&i.Quantity,
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

const getCarnet = `-- name: GetCarnet :one
SELECT id, date, quantity, kid_id FROM carnets
WHERE id=$1
LIMIT 1
`

func (q *Queries) GetCarnet(ctx context.Context, id int64) (Carnet, error) {
	row := q.db.QueryRowContext(ctx, getCarnet, id)
	var i Carnet
	err := row.Scan(
		&i.ID,
		&i.Date,
		&i.Quantity,
		&i.KidID,
	)
	return i, err
}

const updateCarnet = `-- name: UpdateCarnet :one
UPDATE carnets
SET date=$2,
    quantity=$3
WHERE id = $1
RETURNING id, date, quantity, kid_id
`

type UpdateCarnetParams struct {
	ID       int64     `json:"id"`
	Date     time.Time `json:"date"`
	Quantity int32     `json:"quantity"`
}

func (q *Queries) UpdateCarnet(ctx context.Context, arg UpdateCarnetParams) (Carnet, error) {
	row := q.db.QueryRowContext(ctx, updateCarnet, arg.ID, arg.Date, arg.Quantity)
	var i Carnet
	err := row.Scan(
		&i.ID,
		&i.Date,
		&i.Quantity,
		&i.KidID,
	)
	return i, err
}