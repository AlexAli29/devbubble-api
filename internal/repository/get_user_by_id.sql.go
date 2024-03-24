// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: get_user_by_id.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getUserById = `-- name: GetUserById :one
SELECT id, created_at, email, name, is_verified FROM users WHERE id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
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

const getUserIsVerifiedEmailNameById = `-- name: GetUserIsVerifiedEmailNameById :one
SELECT is_verified,email,name FROM users WHERE id = $1
`

type GetUserIsVerifiedEmailNameByIdRow struct {
	IsVerified pgtype.Bool
	Email      string
	Name       string
}

func (q *Queries) GetUserIsVerifiedEmailNameById(ctx context.Context, id pgtype.UUID) (GetUserIsVerifiedEmailNameByIdRow, error) {
	row := q.db.QueryRow(ctx, getUserIsVerifiedEmailNameById, id)
	var i GetUserIsVerifiedEmailNameByIdRow
	err := row.Scan(&i.IsVerified, &i.Email, &i.Name)
	return i, err
}
