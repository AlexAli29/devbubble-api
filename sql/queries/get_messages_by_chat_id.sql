-- name: GetMessagesByChatID :many
SELECT 
    "m"."id" AS "id",
    "m"."text",
    "m"."createdAt",
    "m"."userId",
    "u"."name" AS "senderName",
    ("m"."userId" = $1) AS "isFromMe"
FROM 
    "messages" "m"
JOIN 
    "users" "u" ON "m"."userId" = "u"."id"
WHERE 
    "m"."chatId" = $2
ORDER BY 
    "m"."createdAt" ASC;