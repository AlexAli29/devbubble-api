-- name: RemovePrivateChat :exec
DELETE FROM "private_chats" WHERE "id" = $1;