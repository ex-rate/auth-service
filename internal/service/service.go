package service

import schema "github.com/ex-rate/auth-service/internal/schemas"

type Service struct {
	Registration registration
}

type registration interface {
	CreateUser(reg schema.Registration) error
}
