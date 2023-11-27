package service

import (
	"context"

	"github.com/ex-rate/auth-service/internal/entities"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	passw "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/google/uuid"
)

type service struct {
	user  user
	token token
}

type user interface {
	RegisterUser(ctx context.Context, user schema.Registration) (*schema.Token, error)
	GetUserID(ctx context.Context, username string) (uuid.UUID, error)
}

type token interface {
	CheckToken(ctx context.Context, token string) (string, error)
	GenerateToken(ctx context.Context, user entities.Token) (*schema.Token, error)
}

func New(user user, token token) *service {
	return &service{user, token}
}

func (s *service) RegisterUser(ctx context.Context, user schema.Registration) (*schema.Token, error) {
	password, err := passw.HashPassword(user.HashedPassword)
	if err != nil {
		return nil, err
	}

	user.HashedPassword = password

	return s.user.RegisterUser(ctx, user)
}

func (s *service) RestoreToken(ctx context.Context, token string) (*schema.Token, error) {
	username, err := s.token.CheckToken(ctx, token)
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

	return s.token.GenerateToken(ctx, user)
}
