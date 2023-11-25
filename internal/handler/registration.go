package handler

import (
	"net/http"

	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/gin-gonic/gin"
)

func (h *handler) Registration(ctx *gin.Context) {
	var user schema.Registration
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON", "err": err.Error()})
		return
	}

	// TODO: добавить отправку письма / смс кода на тлф

	ctx.Redirect(http.StatusPermanentRedirect, "/confirm")
}
