package service

import schema "github.com/ex-rate/auth-service/internal/schemas"

func (s *Service) RegisterUser(reg schema.Registration) (string, error) {
	err := s.registration.CreateUser(reg) // создаем юзера в базе
	if err != nil {
		return "", err
	}

	return s.GenerateToken(reg)
}
