package handler

import (
	"context"

	schema "github.com/ex-rate/auth-service/internal/schemas"
)

type handler struct {
	registration registration
}

type registration interface {
	RegisterUser(ctx context.Context, user schema.Registration) (string, error)
}

func New(registration registration) *handler {
	return &handler{registration}
}
