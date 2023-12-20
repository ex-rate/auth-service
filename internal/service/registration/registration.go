package service

import (
	"context"

	"github.com/ex-rate/auth-service/internal/entities"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	token "github.com/ex-rate/auth-service/internal/service/token"

	"github.com/google/uuid"
)

// Registration отвечает за регистрацию пользователей
type Registration struct {
	registrationRepo registrationRepo
	token            *token.Token
}

//go:generate mockgen -source registration.go -destination ../../mocks/registration_repo.go
type registrationRepo interface {
	CreateUser(ctx context.Context, reg schema.Registration) (uuid.UUID, error)
	GetUserID(ctx context.Context, username string) (uuid.UUID, error)
}

func New(registrationRepo registrationRepo, token *token.Token) *Registration {
	return &Registration{registrationRepo: registrationRepo, token: token}
}

// RegisterUser проводит регистрацию пользователя.
// Создает пользователя в базе, возвращает токен
func (s *Registration) RegisterUser(ctx context.Context, user schema.Registration) (*schema.Token, error) {
	userID, err := s.registrationRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	entity := entities.Token{
		UserID:   userID,
		Username: user.Username,
	}

	return s.token.GenerateTokens(ctx, entity)
}

func (s *Registration) GetUserID(ctx context.Context, username string) (uuid.UUID, error) {
	return s.registrationRepo.GetUserID(ctx, username)
}
