package service

import (
	"context"

	"github.com/ex-rate/auth-service/internal/entities"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	registration "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
)

type Service struct {
	user  *registration.Registration
	token *token.Token
}

func New(user *registration.Registration, token *token.Token) *Service {
	return &Service{user, token}
}

// RegisterUser проводит регистрацию пользователя
func (s *Service) RegisterUser(ctx context.Context, user schema.Registration) (*schema.Token, error) {
	password, err := token.HashPassword(user.HashedPassword)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = password

	return s.user.RegisterUser(ctx, user)
}

// RestoreToken проверяет на валидность токен и выдает новый
func (s *Service) RestoreToken(ctx context.Context, token entities.RestoreToken) (*schema.Token, error) {
	err := s.token.CheckAccessToken(token.AccessToken)
	if err != nil {
		return nil, err
	}

	username, err := s.token.CheckRefreshToken(ctx, token.RefreshToken)
	if err != nil {
		return nil, err
	}

	userID, err := s.user.GetUserID(ctx, username)
	if err != nil {
		return nil, err
	}

	user := entities.Token{
		UserID:   userID,
		Username: username,
	}

	return s.token.GenerateTokens(ctx, user)
}
