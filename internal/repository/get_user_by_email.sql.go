// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: get_user_by_email.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, "createdAt", email, description, name, "isVerified" FROM "users" WHERE "email" = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Email,
		&i.Description,
		&i.Name,
		&i.IsVerified,
	)
	return i, err
}

const getUserIsVerifiedEmailNameByEmail = `-- name: GetUserIsVerifiedEmailNameByEmail :one
SELECT "isVerified", "email", "name" FROM "users" WHERE "email" = $1
`

type GetUserIsVerifiedEmailNameByEmailRow struct {
	IsVerified pgtype.Bool `json:"isVerified"`
	Email      string      `json:"email"`
	Name       string      `json:"name"`
}

func (q *Queries) GetUserIsVerifiedEmailNameByEmail(ctx context.Context, email string) (GetUserIsVerifiedEmailNameByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserIsVerifiedEmailNameByEmail, email)
	var i GetUserIsVerifiedEmailNameByEmailRow
	err := row.Scan(&i.IsVerified, &i.Email, &i.Name)
	return i, err
}
