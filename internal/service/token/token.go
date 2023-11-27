package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ex-rate/auth-service/internal/entities"
	api_errors "github.com/ex-rate/auth-service/internal/errors"
	"github.com/golang-jwt/jwt"
)

type token struct {
	secretKey string
	tokenRepo tokenRepo
}

type tokenRepo interface {
	CreateToken(ctx context.Context, token *entities.Token) error
	CheckToken(ctx context.Context, token *entities.Token) error
}

func New(secretKey string, tokenRepo tokenRepo) *token {
	return &token{secretKey: secretKey, tokenRepo: tokenRepo}
}

func (s *token) CheckToken(ctx context.Context, token string) (string, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte("secret"), nil
	})

	if !t.Valid {
		return "", api_errors.ErrInvalidToken
	}

	if err != nil {
		return "", err
	}

	mapClaims := t.Claims.(jwt.MapClaims)
	expr := mapClaims["expr"].(time.Time).Unix()

	entity := &entities.Token{
		RefreshToken: token,
		ExpTime:      expr,
	}

	username := mapClaims["user"].(string)

	return username, s.tokenRepo.CheckToken(ctx, entity)
}
