package handler

import (
	"net/http"

	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/gin-gonic/gin"
)

func (h *handler) Registration(ctx *gin.Context) {
	var reg schema.Registration
	if err := ctx.Bind(&reg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON", "err": err})
	}

	// TODO: добавить отправку письма / смс кода на тлф

	// token, err := h.service.RegisterUser(reg)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"message": "an error occured while creating user", "err": err})
	// }

	//ctx.JSON(http.StatusOK, gin.H{"message": "user successfully created", "token": token})

	ctx.Redirect(http.StatusPermanentRedirect, "/confirm")
}
