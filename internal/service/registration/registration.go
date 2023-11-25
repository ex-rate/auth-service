package service

import (
	"context"

	schema "github.com/ex-rate/auth-service/internal/schemas"
	service "github.com/ex-rate/auth-service/internal/service/token"
)

type registration struct {
	registrationRepo registrationRepo
	token            token
}

//go:generate mockgen -source registration.go -destination ../../mocks/registration_repo.go
type registrationRepo interface {
	CreateUser(ctx context.Context, reg schema.Registration) error
}

type token interface {
	GenerateToken(reg schema.Registration) (string, error)
}

func New(registrationRepo registrationRepo, token token) *registration {
	return &registration{registrationRepo: registrationRepo, token: token}
}

func (s *registration) RegisterUser(ctx context.Context, user schema.Registration) (string, error) {
	password, err := service.HashPassword(user.HashedPassword)
	if err != nil {
		return "", err
	}

	user.HashedPassword = password

	err = s.registrationRepo.CreateUser(ctx, user) // создаем юзера в базе
	if err != nil {
		return "", err
	}

	return s.token.GenerateToken(user)
}
