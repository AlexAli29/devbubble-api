-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserIsVerifiedEmailNameById :one
SELECT is_verified,email,name FROM users WHERE id = $1;