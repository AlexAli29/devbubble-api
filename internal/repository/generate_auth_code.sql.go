// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: generate_auth_code.sql

package db

import (
	"context"
)

const generateAuthCode = `-- name: GenerateAuthCode :one
WITH updated AS (
  UPDATE "auth_codes"
  SET "code" = floor(random() * 900000 + 100000)     
  FROM "users"  
  WHERE "users"."email" = $1 
    AND "auth_codes"."userId" = "users"."id" 
  RETURNING "auth_codes"."code"
)
SELECT "code" FROM updated
`

func (q *Queries) GenerateAuthCode(ctx context.Context, email string) (int32, error) {
	row := q.db.QueryRow(ctx, generateAuthCode, email)
	var code int32
	err := row.Scan(&code)
	return code, err
}
