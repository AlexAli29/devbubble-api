// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: create_auth_code.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAuthCode = `-- name: CreateAuthCode :one
INSERT INTO auth_codes (user_id)
VALUES ($1)
RETURNING code
`

func (q *Queries) CreateAuthCode(ctx context.Context, userID pgtype.UUID) (int32, error) {
	row := q.db.QueryRow(ctx, createAuthCode, userID)
	var code int32
	err := row.Scan(&code)
	return code, err
}