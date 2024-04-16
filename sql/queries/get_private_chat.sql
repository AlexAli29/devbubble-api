-- name: GetPrivateChat :one
SELECT 
    "pc"."chatId" AS "chatId",
    "u"."id" AS "otherParticipantId",
    "u"."name" AS "otherParticipantName"
FROM 
    "chat_participants" "pc"
JOIN 
    "users" "u" ON "u"."id" = "pc"."userId"
WHERE 
    "pc"."chatId" = $2
    AND "pc"."userId" <> $1;