package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ex-rate/auth-service/internal/entities"
	api_errors "github.com/ex-rate/auth-service/internal/errors"
	"github.com/gin-gonic/gin"
)

const AuthorizationHeader = "Authorization"

func (h *handler) RestoreToken(ctx *gin.Context) {
	var token entities.RestoreToken

	// забираем акссес токен из хедера
	accessTokenString := ctx.Request.Header.Get(AuthorizationHeader)
	accessToken := strings.Split(accessTokenString, "Bearer ")[1]
	token.AccessToken = accessToken

	// забираем рефреш токен из тела запроса
	dec := json.NewDecoder(ctx.Request.Body)
	err := dec.Decode(&token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	newToken, err := h.service.RestoreToken(ctx.Request.Context(), token)
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

	jsonMsg := gin.H{"message": "successfully created token", "access-token": newToken.AccessToken, "refresh-token": newToken.RefreshToken}

	ctx.JSON(http.StatusOK, jsonMsg)
}
