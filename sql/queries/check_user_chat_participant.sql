-- name: IsUserChatParticipant :one
SELECT 
    EXISTS (
        SELECT 1 
        FROM "chat_participants" 
        WHERE "chatId" = $1 AND "userId" = $2
    ) AS "is_participant";