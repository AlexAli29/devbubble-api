-- name: CreateAuthCode :one
INSERT INTO auth_codes ("userId")
VALUES ($1)
RETURNING "code";
