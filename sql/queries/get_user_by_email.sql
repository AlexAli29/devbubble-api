-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserIsVerifiedEmailNameByEmail :one
SELECT is_verified,email,name FROM users WHERE email = $1;