package customRepository

import (
	"context"
	"devbubble-api/internal/core"
	"devbubble-api/pkg/sql"

	"github.com/jackc/pgx/v5/pgtype"
)

const getPrivateChats = `SELECT 
pc.id AS "chatId",
u.name AS "participantName",
m.text AS "lastMessage",
m."createdAt" AS "lastMessageTime",
(m.user_id = $1) AS "isFromMe",
sender.name AS "lastMessageSenderName"
FROM 
chat_participants cp
JOIN 
private_chats pc ON cp."chatId" = pc.id
LEFT JOIN 
users u ON cp."userId" = u.id
LEFT JOIN 
(
    SELECT 
        "chatId",
        "userId",
        text,
        "createdAt"
    FROM 
        messages
    WHERE 
        "createdAt" IN (
            SELECT MAX("createdAt")
            FROM messages
            GROUP BY "chatId"
        )
) m ON cp."chatId" = m."chatId"
LEFT JOIN 
users sender ON sender.id = m."userId"  -- Дополнительное соединение для получения имени отправителя
WHERE 
cp."userId" <> $1
AND pc.id IN (
    SELECT "chatId" 
    FROM chat_participants 
    WHERE "userId" = $1
)
ORDER BY 
m.createdAt DESC NULLS LAST;
`

type PrivateChat struct {
	ID                    pgtype.UUID      `json:"chatId"`
	ParticipantName       string           `json:"participantName"`
	LastMessage           pgtype.Text      `json:"lastMessage"`
	LastMessageTime       pgtype.Timestamp `json:"lastMessageTime"`
	IsFromMe              pgtype.Bool      `json:"isFromMe"`
	LastMessageSenderName pgtype.Text      `json:"lastMessageSenderName"`
}

func (q *Queries) GetPrivateChats(ctx context.Context, userId pgtype.UUID) ([]core.GetPrivateChatsResponse, error) {
	rows, err := q.db.Query(ctx, getPrivateChats, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []core.GetPrivateChatsResponse

	for rows.Next() {
		var i PrivateChat
		if err := rows.Scan(&i.ID, &i.ParticipantName, &i.LastMessage, &i.LastMessageTime, &i.IsFromMe, &i.LastMessageSenderName); err != nil {
			return nil, err
		}
		items = append(items, core.GetPrivateChatsResponse{Id: sql.UUIDToString(i.ID),
			Name:            i.ParticipantName,
			LastMessage:     i.LastMessage.String,
			LastMessageTime: i.LastMessageTime, IsFromMe: i.IsFromMe.Bool, LastMessageSenderName: i.LastMessageSenderName, ChatType: core.PRIVATE_CHAT})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
