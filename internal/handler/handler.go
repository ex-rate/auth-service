package handler

import (
	"context"

	schema "github.com/ex-rate/auth-service/internal/schemas"
)

type handler struct {
	service srv
}

type srv interface {
	RegisterUser(ctx context.Context, user schema.Registration) (*schema.Token, error)
	RestoreToken(ctx context.Context, token string) (*schema.Token, error)
}

func New(service srv) *handler {
	return &handler{service}
}
