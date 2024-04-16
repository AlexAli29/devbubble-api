-- name: AddPrivateChatParticipants :exec
INSERT INTO chat_participants ("chatId", "userId")
VALUES
    ($1, $2),
    ($1, $3);
