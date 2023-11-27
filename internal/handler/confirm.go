package handler

import (
	"errors"
	"net/http"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/gin-gonic/gin"
)

func (h *handler) Confirm(ctx *gin.Context) {
	var user schema.Registration
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "invalid JSON", "err": err.Error()})
		return
	}

	// проверка кода

	c := ctx.Request.Context()

	token, err := h.service.RegisterUser(c, user)
	if err != nil {
		if errors.Is(err, api_errors.ErrUsernameAlreadyExists) || errors.Is(err, api_errors.ErrEmailAlreadyExists) ||
			errors.Is(err, api_errors.ErrPhoneAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "an error occured while creating user", "err": err.Error()})
		return
	}

	jsonMsg := gin.H{"message": "user successfully created", "access-token": token.AccessToken, "refresh-token": token.RefreshToken}

	ctx.JSON(http.StatusOK, jsonMsg)
}
