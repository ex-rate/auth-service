package service

import schema "github.com/ex-rate/auth-service/internal/schemas"

type Service struct {
	registration registration
	secretKey    string
}

type registration interface {
	CreateUser(reg schema.Registration) error
}

func New(registration registration, secretKey string) *Service {
	return &Service{registration: registration, secretKey: secretKey}
}
