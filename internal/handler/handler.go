package handler

import (
	"github.com/ex-rate/auth-service/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service}
}
