package router

import (
	"github.com/ex-rate/auth-service/internal/handler"
	"github.com/gin-gonic/gin"
)

// New создает роутер и настраивает пути
func New(handler *handler.Handler) *gin.Engine {
	r := gin.Default()

	// registration
	r.GET("/signup", handler.GetRegistration)
	r.POST("/signup", handler.Registration)

	// confirm registration
	r.GET("/confirm", handler.GetConfirm)
	r.POST("/confirm", handler.Confirm)

	// restore token
	r.PUT("/restore_token", handler.RestoreToken)

	// authorization
	r.GET("/login", handler.GetAuth)
	r.POST("/login", handler.Auth)

	// authorization via code
	r.GET("/code", handler.GetCode)
	r.POST("/code", handler.Code)

	// authorization via password
	r.GET("/password", handler.GetPassword)
	r.POST("/password", handler.Password)

	return r
}
