-- name: CreateMessage :one
INSERT INTO messages ("userId", "chatId", text)
VALUES ($1, $2, $3)
RETURNING *;