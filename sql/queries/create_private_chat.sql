-- name: CreatePrivateChat :one
WITH new_chat AS (
    INSERT INTO "private_chats" DEFAULT VALUES 
    RETURNING "id"
)
SELECT "id" FROM new_chat;