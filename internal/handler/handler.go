package handler

import "github.com/ex-rate/auth-service/internal/service"

type handler struct {
	service *service.Service
}

func New(service *service.Service) *handler {
	return &handler{service: service}
}
