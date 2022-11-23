// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: teacher.sql

package db

import (
	"context"
)

const createTeacher = `-- name: CreateTeacher :one
INSERT INTO teachers(name, surname)
VALUES ($1, $2)
RETURNING id, name, surname
`

type CreateTeacherParams struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

func (q *Queries) CreateTeacher(ctx context.Context, arg CreateTeacherParams) (Teacher, error) {
	row := q.db.QueryRowContext(ctx, createTeacher, arg.Name, arg.Surname)
	var i Teacher
	err := row.Scan(&i.ID, &i.Name, &i.Surname)
	return i, err
}

const getAllTeachers = `-- name: GetAllTeachers :many
SELECT id, name, surname FROM teachers
ORDER BY name
`

func (q *Queries) GetAllTeachers(ctx context.Context) ([]Teacher, error) {
	rows, err := q.db.QueryContext(ctx, getAllTeachers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Teacher{}
	for rows.Next() {
		var i Teacher
		if err := rows.Scan(&i.ID, &i.Name, &i.Surname); err != nil {
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

const getTeacher = `-- name: GetTeacher :one
SELECT id, name, surname FROM teachers
WHERE id=$1
LIMIT 1
`

func (q *Queries) GetTeacher(ctx context.Context, id int64) (Teacher, error) {
	row := q.db.QueryRowContext(ctx, getTeacher, id)
	var i Teacher
	err := row.Scan(&i.ID, &i.Name, &i.Surname)
	return i, err
}
