package service

import (
	"context"
	"devbubble-api/internal/config"
	"devbubble-api/internal/core"
	db "devbubble-api/internal/repository"
	"devbubble-api/pkg/sql"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
)

type MessageService struct {
	queries *db.Queries
	log     *slog.Logger
	cfg     *config.Config
}

func NewMessageService(queries *db.Queries, log *slog.Logger, cfg *config.Config) *MessageService {

	return &MessageService{queries, log, cfg}

}

func (s *MessageService) CreateMessage(ctx context.Context, dto core.CreateMessageDto) (core.Message, error) {
	userUUID := pgtype.UUID{}
	chatUUID := pgtype.UUID{}
	err := userUUID.Scan(dto.UserId)
	if err != nil {
		return core.Message{}, err
	}
	err = chatUUID.Scan(dto.ChatId)
	if err != nil {
		return core.Message{}, err
	}
	user, err := s.queries.GetUserIsVerifiedEmailNameById(ctx, userUUID)
	if err != nil {
		return core.Message{}, err
	}

	newMessage, err := s.queries.CreateMessage(ctx, db.CreateMessageParams{UserId: userUUID, ChatId: chatUUID, Text: dto.Text})
	if err != nil {
		return core.Message{}, err
	}
	return core.Message{
		Id:         sql.UUIDToString(newMessage.ID),
		ChatId:     dto.ChatId,
		SenderId:   dto.UserId,
		SenderName: user.Name,
		IsFromMe:   false,
		Text:       newMessage.Text,
		CreatedAt:  newMessage.CreatedAt.Time.String(),
	}, nil
}

func (s *MessageService) RemoveMessage(ctx context.Context, dto core.RemoveMessageDto) error {
	userUUID := pgtype.UUID{}
	messageUUID := pgtype.UUID{}

	err := userUUID.Scan(dto.UserId)
	if err != nil {
		return err
	}
	err = messageUUID.Scan(dto.Id)
	if err != nil {
		return err
	}
	err = s.queries.RemoveMessage(ctx, db.RemoveMessageParams{ID: messageUUID, UserId: userUUID})
	if err != nil {
		return err
	}
	return nil
}

func (s *MessageService) GetMessagesByChatId(ctx context.Context, dto core.GetMessagesByChatIdDto) ([]db.GetMessagesByChatIDRow, error) {
	userUUID := pgtype.UUID{}
	chatUUID := pgtype.UUID{}

	err := userUUID.Scan(dto.UserId)
	if err != nil {
		return nil, err
	}
	err = chatUUID.Scan(dto.ChatId)
	if err != nil {
		return nil, err
	}
	messages, err := s.queries.GetMessagesByChatID(ctx, db.GetMessagesByChatIDParams{ChatId: chatUUID, UserId: userUUID})
	if err != nil {
		return nil, err
	}
	return messages, nil
}
