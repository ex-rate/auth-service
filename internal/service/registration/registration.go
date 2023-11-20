package service

import schema "github.com/ex-rate/auth-service/internal/schemas"

type registration struct {
	registrationRepo registrationRepo
	token            token
}

type registrationRepo interface {
	CreateUser(reg schema.Registration) error
}

type token interface {
	GenerateToken(reg schema.Registration) (string, error)
}

func New(registrationRepo registrationRepo, token token) *registration {
	return &registration{registrationRepo: registrationRepo, token: token}
}

func (s *registration) RegisterUser(user schema.Registration) (string, error) {
	err := s.registrationRepo.CreateUser(user) // создаем юзера в базе
	if err != nil {
		return "", err
	}

	return s.token.GenerateToken(user)
}
