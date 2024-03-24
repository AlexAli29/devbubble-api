// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: create_user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1,$2)
RETURNING id, created_at, email, name, is_verified
`

type CreateUserParams struct {
	Name  string
	Email string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Name, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Email,
		&i.Name,
		&i.IsVerified,
	)
	return i, err
}