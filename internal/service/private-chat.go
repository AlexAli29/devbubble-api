package service

import (
	"context"
	"devbubble-api/internal/config"
	"devbubble-api/internal/core"
	customRepository "devbubble-api/internal/custom-repository"
	db "devbubble-api/internal/repository"
	"devbubble-api/pkg/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PrivateChatService struct {
	queries *db.Queries

	log           *slog.Logger
	cfg           *config.Config
	pool          *pgxpool.Pool
	customQueries *customRepository.Queries
}

func NewPrivateChatService(queries *db.Queries, log *slog.Logger, cfg *config.Config, pool *pgxpool.Pool, customQueries *customRepository.Queries) *PrivateChatService {

	return &PrivateChatService{queries, log, cfg, pool, customQueries}

}

func (s *PrivateChatService) CreatePrivateChat(ctx context.Context, dto core.CreatePrivateChatDto) (string, error) {
	tx, err := s.pool.Begin(ctx)
	qtx := s.queries.WithTx(tx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	chatId, err := qtx.CreatePrivateChat(ctx)
	if err != nil {
		return "", err
	}

	userUUID := pgtype.UUID{}
	SecondParticipantUUID := pgtype.UUID{}

	err = userUUID.Scan(dto.UserId)
	if err != nil {
		return "", err
	}
	err = SecondParticipantUUID.Scan(dto.SecondParticipantId)
	if err != nil {
		return "", err
	}

	err = qtx.AddPrivateChatParticipants(ctx, db.AddPrivateChatParticipantsParams{ChatId: chatId, UserId: userUUID, UserId_2: SecondParticipantUUID})
	if err != nil {
		return "", err
	}

	tx.Commit(ctx)

	return sql.UUIDToString(chatId), nil
}
func (s *PrivateChatService) RemovePrivateChat(ctx context.Context, dto core.RemovePrivateChatDto) error {
	chatUUID := pgtype.UUID{}
	userUUID := pgtype.UUID{}
	err := chatUUID.Scan(dto.ChatId)
	if err != nil {
		return err
	}
	err = userUUID.Scan(dto.UserId)
	if err != nil {
		return err
	}
	isUserChatParticipant, err := s.queries.IsUserChatParticipant(ctx, db.IsUserChatParticipantParams{ChatId: chatUUID, UserId: userUUID})
	if err != nil {
		return err
	}
	if !isUserChatParticipant {
		return errors.New("user must be chat participant to remove chat")
	}
	err = s.queries.RemovePrivateChat(ctx, chatUUID)
	if err != nil {
		return err
	}
	return nil
}

func (s *PrivateChatService) GetPrivateChats(ctx context.Context, userId string) ([]core.GetPrivateChatsResponse, error) {

	userUUID := pgtype.UUID{}
	err := userUUID.Scan(userId)
	if err != nil {
		return []core.GetPrivateChatsResponse{}, err
	}

	chats, err := s.customQueries.GetPrivateChats(ctx, userUUID)
	fmt.Printf("uu %v", chats)
	if err != nil {
		return []core.GetPrivateChatsResponse{}, err
	}
	return chats, nil
}
