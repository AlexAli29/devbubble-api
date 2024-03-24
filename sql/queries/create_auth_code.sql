-- name: CreateAuthCode :one
INSERT INTO auth_codes (user_id)
VALUES ($1)
RETURNING code;
