package service

import (
	"context"
	"devbubble-api/internal/config"
	"devbubble-api/internal/core"
	db "devbubble-api/internal/repository"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
)

type TagService struct {
	queries *db.Queries

	log *slog.Logger
	cfg *config.Config
}

func NewTagService(queries *db.Queries, log *slog.Logger, cfg *config.Config) *TagService {

	return &TagService{queries, log, cfg}

}

func (s *TagService) GetAvailableTags(ctx context.Context, userId string) ([]db.UserTag, error) {
	userIdUUID := pgtype.UUID{}

	err := userIdUUID.Scan(userId)

	if err != nil {
		return []db.UserTag{}, err
	}
	tagsDb, err := s.queries.GetAvailableTags(ctx, userIdUUID)
	if err != nil {
		return []db.UserTag{}, err
	}
	return tagsDb, nil
}

func (s *TagService) AddTag(ctx context.Context, dto core.AddTagDto) error {
	userIdUUID := pgtype.UUID{}
	tagIdUUID := pgtype.UUID{}
	err := userIdUUID.Scan(dto.UserId)
	if err != nil {
		return err
	}
	s.log.Debug(fmt.Sprintf("%v", dto))
	err = tagIdUUID.Scan(dto.TagId)
	if err != nil {
		return err
	}
	err = s.queries.AddUserTag(ctx, db.AddUserTagParams{UserId: userIdUUID, UserTagId: tagIdUUID})
	if err != nil {
		return err
	}
	return nil
}
func (s *TagService) RemoveTag(ctx context.Context, dto core.RemoveTagDto) error {
	userIdUUID := pgtype.UUID{}
	tagIdUUID := pgtype.UUID{}
	err := userIdUUID.Scan(dto.UserId)
	if err != nil {
		return err
	}
	s.log.Debug(fmt.Sprintf("%v", dto))
	err = tagIdUUID.Scan(dto.TagId)
	if err != nil {
		return err
	}
	err = s.queries.DeleteUserTag(ctx, db.DeleteUserTagParams{UserId: userIdUUID, UserTagId: tagIdUUID})
	if err != nil {
		return err
	}
	return nil
}
