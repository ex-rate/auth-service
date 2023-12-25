package service

import (
	"context"
	"fmt"

	"github.com/ex-rate/auth-service/internal/entities"
	api_errors "github.com/ex-rate/auth-service/internal/errors"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// Token - структура, создающая и проверяющая токен
type Token struct {
	secretKey string
	tokenRepo tokenRepo
}

//go:generate mockgen -source token.go -destination ../../mocks/token_repo.go
type tokenRepo interface {
	CreateToken(ctx context.Context, token *entities.Token) error
	CheckToken(ctx context.Context, token *entities.Token) error
	GetUserID(ctx context.Context, username string) (uuid.UUID, error)
}

func New(secretKey string, tokenRepo tokenRepo) *Token {
	return &Token{secretKey: secretKey, tokenRepo: tokenRepo}
}

// checkToken проверяет токен на валидность; возвращает username
func (s *Token) checkToken(token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(s.secretKey), nil
	})

	if !t.Valid {
		return "", api_errors.ErrInvalidToken
	}

	if err != nil {
		return "", err
	}

	mapClaims := t.Claims.(jwt.MapClaims)

	username := mapClaims["user"].(string)

	isAuthorized := mapClaims["authorized"].(bool)
	if !isAuthorized {
		return "", api_errors.ErrNotAuthorized
	}

	return username, nil
}

// CheckTokens проверяет на валидность два токена: access & refresh
func (s *Token) CheckTokens(ctx context.Context, token schema.RestoreToken) (string, error) {
	accessUsername, err := s.checkToken(token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("error while checking access token: %w", err)
	}

	refreshUsername, err := s.checkToken(token.RefreshToken)
	if err != nil {
		return "", fmt.Errorf("error while checking refresh token: %w", err)
	}

	if refreshUsername != accessUsername {
		return "", api_errors.ErrInvalidUsername
	}

	t, err := jwt.Parse(token.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	mapClaims := t.Claims.(jwt.MapClaims)
	expr := mapClaims["exp"].(float64)

	entity := &entities.Token{
		RefreshToken: token.RefreshToken,
		ExpTime:      expr,
	}

	return refreshUsername, s.tokenRepo.CheckToken(ctx, entity)
}
