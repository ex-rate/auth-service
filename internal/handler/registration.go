package handler

import (
	"net/http"

	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/gin-gonic/gin"
)

// GET /signup
func (h *Handler) GetRegistration(ctx *gin.Context) {
	// TODO: html render
}

// POST /signup
func (h *Handler) Registration(ctx *gin.Context) {
	var user schema.Registration
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// TODO: добавить отправку письма / смс кода на тлф

	ctx.Redirect(http.StatusMovedPermanently, "/confirm")
}
