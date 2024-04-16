-- name: AddUserTag :exec
INSERT INTO "user_user_tags" ("userId", "userTagId")
VALUES ($1, $2)
ON CONFLICT ("userId", "userTagId") DO NOTHING;

-- name: DeleteUserTag :exec
DELETE FROM "user_user_tags"
WHERE "userId" = $1 AND "userTagId" = $2;

-- name: UpdateUser :exec
UPDATE users
SET
  name = COALESCE(NULLIF(sqlc.arg(name)::text, ''), name),
  description = COALESCE(NULLIF(sqlc.arg(description)::text, ''), description)
WHERE id = $1;