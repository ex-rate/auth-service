package service

import (
	"context"
	"fmt"

	"github.com/ex-rate/auth-service/internal/entities"
	api_errors "github.com/ex-rate/auth-service/internal/errors"
	"github.com/golang-jwt/jwt"
)

// token - структура, создающая и проверяющая токен
type token struct {
	secretKey string
	tokenRepo tokenRepo
}

//go:generate mockgen -source token.go -destination ../../mocks/token_repo.go
type tokenRepo interface {
	CreateToken(ctx context.Context, token *entities.Token) error
	CheckToken(ctx context.Context, token *entities.Token) error
}

func New(secretKey string, tokenRepo tokenRepo) *token {
	return &token{secretKey: secretKey, tokenRepo: tokenRepo}
}

// CheckToken проверяет токен на валидность.
// Возвращает username пользователя, которому принадлежит токен и ошибку
func (s *token) CheckToken(ctx context.Context, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !t.Valid {
		return "", api_errors.ErrInvalidToken
	}

	mapClaims := t.Claims.(jwt.MapClaims)
	expr := mapClaims["exp"].(float64)

	entity := &entities.Token{
		RefreshToken: token,
		ExpTime:      expr,
	}

	username := mapClaims["user"].(string)

	return username, s.tokenRepo.CheckToken(ctx, entity)
}
