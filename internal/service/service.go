package service

import (
	"context"

	"github.com/ex-rate/auth-service/internal/entities"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/ex-rate/auth-service/internal/service/auth"
	registration "github.com/ex-rate/auth-service/internal/service/registration"
	token "github.com/ex-rate/auth-service/internal/service/token"
)

type Service struct {
	user  *registration.Registration
	token *token.Token
	auth  *auth.AuthService
}

func New(user *registration.Registration, token *token.Token, auth *auth.AuthService) *Service {
	return &Service{user, token, auth}
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
func (s *Service) RestoreToken(ctx context.Context, token schema.RestoreToken) (*schema.Token, error) {
	username, err := s.token.CheckTokens(ctx, token)
	if err != nil {
		return nil, err
	}

	user := entities.Token{
		Username: username,
	}

	return s.token.GenerateTokens(ctx, user)
}

// AuthWithCode проводит авторизацию пользователя по коду подтверждения через смс / почту
func (s *Service) AuthWithCode(ctx context.Context, user schema.AuthWithCode) (*schema.Token, error) {
	return s.auth.WithCode(ctx, user)
}

// AuthWithPassword проводит авторизацию пользователя по паролю
func (s *Service) AuthWithPassword(ctx context.Context, user schema.AuthWithPassword) (*schema.Token, error) {
	return s.auth.WithPassword(ctx, user)
}
