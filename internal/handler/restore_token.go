package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/gin-gonic/gin"
)

const AuthorizationHeader = "Authorization"

func (h *handler) RestoreToken(ctx *gin.Context) {
	var token schema.RestoreToken

	// забираем акссес токен из хедера
	accessTokenString := ctx.Request.Header.Get(AuthorizationHeader)
	accessToken := strings.Split(accessTokenString, "Bearer ")[1]

	// забираем рефреш токен из тела запроса
	dec := json.NewDecoder(ctx.Request.Body)
	err := dec.Decode(&token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	token.AccessToken = accessToken

	newToken, err := h.service.RestoreToken(ctx.Request.Context(), token)
	if err != nil {
		if err != nil {
			if errors.Is(err, api_errors.ErrInvalidToken) || errors.Is(err, api_errors.ErrInvalidUsername) ||
				errors.Is(err, api_errors.ErrTokenNotExists) {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
				return
			}

			if errors.Is(err, api_errors.ErrNotAuthorized) {
				ctx.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	jsonMsg := gin.H{"message": "successfully created token", "access-token": newToken.AccessToken, "refresh-token": newToken.RefreshToken}

	ctx.JSON(http.StatusOK, jsonMsg)
}
