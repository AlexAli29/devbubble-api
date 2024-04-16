-- name: RemoveMessage :exec
DELETE FROM "messages"
WHERE "id" = $1 AND "userId" = $2;