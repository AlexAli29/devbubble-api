// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: create_private_chat.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPrivateChat = `-- name: CreatePrivateChat :one
WITH new_chat AS (
    INSERT INTO "private_chats" DEFAULT VALUES 
    RETURNING "id"
)
SELECT "id" FROM new_chat
`

func (q *Queries) CreatePrivateChat(ctx context.Context) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createPrivateChat)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}
