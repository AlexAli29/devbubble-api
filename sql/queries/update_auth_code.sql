-- name: UpdateAuthCode :one
UPDATE "auth_codes"
SET "code" = $1  -- Replace $1 with the new code value
WHERE "userId" = $2  -- Replace $2 with the user ID
RETURNING "code";  -- Optional: retrieve the updated code