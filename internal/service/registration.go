package service

import schema "github.com/ex-rate/auth-service/internal/schemas"

func (s *Service) RegisterUser(reg schema.Registration) error {
	return s.Registration.CreateUser(reg)
}
