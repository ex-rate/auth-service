package handler

import (
	"errors"
	"net/http"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	"github.com/gin-gonic/gin"
)

func (h *handler) RestoreToken(ctx *gin.Context) {
	refreshToken := ctx.Request.Header.Get("refresh-token")

	token, err := h.service.RestoreToken(ctx.Request.Context(), refreshToken)
	if err != nil {
		if err != nil {
			if errors.Is(err, api_errors.ErrInvalidToken) {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	jsonMsg := gin.H{"message": "user successfully created", "access-token": token.AccessToken, "refresh-token": token.RefreshToken}

	ctx.JSON(http.StatusOK, jsonMsg)
}
