package service

import (
	"context"
	"devbubble-api/internal/core"
	db "devbubble-api/internal/repository"
	"devbubble-api/pkg/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgtype"
)

type UserService struct {
	queries     *db.Queries
	authService *AuthService
	log         *slog.Logger
}

func NewUserService(queries *db.Queries, authService *AuthService, log *slog.Logger) *UserService {

	return &UserService{queries, authService, log}
}

func (s *UserService) UpdateUser(ctx context.Context, userDto core.UpdateUserDto) error {
	userIdUUID := pgtype.UUID{}

	err := userIdUUID.Scan(userDto.Id)
	if err != nil {
		return err
	}
	err = s.queries.UpdateUser(ctx, db.UpdateUserParams{ID: userIdUUID, Name: userDto.Name, Description: userDto.Description})
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) GetUserById(id int) string {
	return "GetUserById"
}
func (s *UserService) CreateUser(ctx context.Context, userDto core.CreateUserDto) (string, error) {
	userByEmail, err := s.queries.GetUserByEmail(ctx, userDto.Email)
	if err == nil && userByEmail.IsVerified.Bool {
		return "", errors.New("email already taken")
	}

	if err == nil && !userByEmail.IsVerified.Bool {
		go func() {
			ctx2 := context.WithoutCancel(ctx)
			s.authService.GenerateAuthCode(ctx2, userByEmail.Email)
		}()
		return userByEmail.Email, nil
	}

	user, err := s.queries.CreateUser(ctx, db.CreateUserParams{Name: userDto.Name, Email: userDto.Email})
	s.log.Debug(fmt.Sprintf("%v", user))
	if err != nil {
		return "", err
	}

	go func() {
		ctx2 := context.WithoutCancel(ctx)
		s.authService.CreateAuthCode(ctx2, core.CreateAuthCodeDto{UserId: sql.UUIDToString(user.ID), Email: user.Email, Name: user.Name})
	}()

	return user.Email, nil

}

func (s *UserService) VerifyUser(ctx context.Context, dto core.VerifyUserDto) (string, error) {
	id, err := s.queries.VerifyUser(ctx, db.VerifyUserParams{Email: dto.Email, Code: dto.Code})
	if err != nil {

		s.log.Error(err.Error())
		return "", err

	}
	stringId := sql.UUIDToString(id)
	return stringId, nil
}

func (s *UserService) GetCurrentUser(ctx context.Context, id string) (core.CurrentUserResponse, error) {
	userIdUUID := pgtype.UUID{}

	err := userIdUUID.Scan(id)

	if err != nil {
		return core.CurrentUserResponse{}, err
	}
	dbUser, err := s.queries.GetUserById(ctx, userIdUUID)
	if err != nil {
		return core.CurrentUserResponse{}, err
	}
	s.log.Debug(fmt.Sprintf("%v", dbUser.Tags))
	var tags []core.Tag

	// Unmarshal the JSON-encoded tags
	if len(dbUser.Tags) > 0 {
		err = json.Unmarshal(dbUser.Tags, &tags)
		if err != nil {
			// Handle error
			fmt.Println("Error unmarshalling tags:", err)
			return core.CurrentUserResponse{}, err
		}
	}
	return core.CurrentUserResponse{Id: sql.UUIDToString(dbUser.ID), Email: dbUser.Email, Name: dbUser.Name, Description: dbUser.Description.String, Tags: tags}, nil
}
