package service

import (
	"context"

	"github.com/ex-rate/auth-service/internal/entities"
	schema "github.com/ex-rate/auth-service/internal/schemas"

	"github.com/google/uuid"
)

// registration отвечает за регистрацию пользователей
type registration struct {
	registrationRepo registrationRepo
	token            token
}

//go:generate mockgen -source registration.go -destination ../../mocks/registration_repo.go
type registrationRepo interface {
	CreateUser(ctx context.Context, reg schema.Registration) (uuid.UUID, error)
	GetUserID(ctx context.Context, username string) (uuid.UUID, error)
}

type token interface {
	GenerateToken(ctx context.Context, user entities.Token) (*schema.Token, error)
}

func New(registrationRepo registrationRepo, token token) *registration {
	return &registration{registrationRepo: registrationRepo, token: token}
}

// RegisterUser проводит регистрацию пользователя.
// Создает пользователя в базе, возвращает токен
func (s *registration) RegisterUser(ctx context.Context, user schema.Registration) (*schema.Token, error) {
	userID, err := s.registrationRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	entity := entities.Token{
		UserID:   userID,
		Username: user.Username,
	}

	return s.token.GenerateToken(ctx, entity)
}

func (s *registration) GetUserID(ctx context.Context, username string) (uuid.UUID, error) {
	return s.registrationRepo.GetUserID(ctx, username)
}
