package service

import (
	"context"
	"devbubble-api/internal/config"
	"devbubble-api/internal/core"
	db "devbubble-api/internal/repository"
	"devbubble-api/pkg/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthService struct {
	queries      *db.Queries
	emailService *EmailService
	log          *slog.Logger
	cfg          *config.Config
}

func NewAuthService(queries *db.Queries, emailService *EmailService, log *slog.Logger, cfg *config.Config) *AuthService {

	return &AuthService{queries, emailService, log, cfg}

}

func (s *AuthService) CreateAuthCode(ctx context.Context, dto core.CreateAuthCodeDto) {

	userIdUUID := pgtype.UUID{}

	err := userIdUUID.Scan(dto.UserId)

	if err != nil {
		return
	}

	code, err := s.queries.CreateAuthCode(ctx, userIdUUID)
	ctx.Done()
	if err != nil {
		s.log.Error(err.Error())
	}

	s.log.Debug(fmt.Sprintf("code : %d", code))
	s.log.Debug(fmt.Sprintf("email : %v", dto))

	s.emailService.SendTextEmail(dto.Email, fmt.Sprintf("Hello %s <br/> This is your auth code: <b>%d</b>", dto.Name, code), "Authorization code")
}

func (s *AuthService) GenerateAuthCode(ctx context.Context, email string) error {
	code, err := s.queries.GenerateAuthCode(ctx, email)
	if err != nil {
		return err
	}
	go s.emailService.SendTextEmail(email, fmt.Sprintf("This is your auth code: <b>%d</b>", code), "Authorization code")
	return nil
}

func (s *AuthService) SignIn(ctx context.Context, dto core.SignInDto) (string, error) {

	id, err := s.queries.CheckAuthCode(ctx, db.CheckAuthCodeParams{Email: dto.Email, Code: dto.Code})
	if err != nil {
		return "", err
	}
	stringId := sql.UUIDToString(id)

	return stringId, nil
}

func (s *AuthService) CreateJWT(userId string) (string, error) {
	claims := core.AuthTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(s.cfg.JWT.ExpirationTimeHours)).Unix(),
			Issuer:    "devBubble",
		},
		UserId: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWT.Secret))

}

func (s *AuthService) ParseToken(tokenString string) (*core.AuthTokenClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(tokenString, &core.AuthTokenClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.cfg.JWT.Secret), nil
	})

	if claims, ok := jwtToken.Claims.(*core.AuthTokenClaims); ok && jwtToken.Valid {

		return claims, nil
	} else {

		return nil, err
	}

}
