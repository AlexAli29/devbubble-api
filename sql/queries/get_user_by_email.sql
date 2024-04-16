-- name: GetUserByEmail :one
SELECT * FROM "users" WHERE "email" = $1;

-- name: GetUserIsVerifiedEmailNameByEmail :one
SELECT "isVerified", "email", "name" FROM "users" WHERE "email" = $1;