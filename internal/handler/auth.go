package handler

import (
	"errors"
	"net/http"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/gin-gonic/gin"
)

// GET /login
func (h *Handler) GetAuth(ctx *gin.Context) {
	// TODO: render html
}

// POST /login
func (h *Handler) Auth(ctx *gin.Context) {
	// var auth schema.Auth

	// if err := ctx.BindJSON(&auth); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	// 	return
	// }

	// редирект на страницу с вводом кода
	ctx.Redirect(http.StatusMovedPermanently, "/code")
}

// GET /code
func (h *Handler) GetCode(ctx *gin.Context) {
	// TODO: render html
}

// POST /code
func (h *Handler) Code(ctx *gin.Context) {
	var auth schema.AuthWithCode

	if err := ctx.BindJSON(&auth); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := h.service.AuthWithCode(ctx.Request.Context(), auth)
	if err != nil {
		if errors.Is(err, api_errors.ErrUserNotExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	jsonMsg := gin.H{
		"message":       "successfully logged in",
		"access-token":  token.AccessToken,
		"refresh-token": token.RefreshToken}

	ctx.JSON(http.StatusOK, jsonMsg)
}

// GetPassword используется, если пользователь не может войти по коду подтверждения.
// GET /password
func (h *Handler) GetPassword(ctx *gin.Context) {
	// TODO: render html
}

// POST /password
func (h *Handler) Password(ctx *gin.Context) {
	var auth schema.AuthWithPassword

	if err := ctx.BindJSON(&auth); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	token, err := h.service.AuthWithPassword(ctx.Request.Context(), auth)
	if err != nil {
		if errors.Is(err, api_errors.ErrUserNotExists) || errors.Is(err, api_errors.ErrIncorrectPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	jsonMsg := gin.H{
		"message":       "successfully logged in",
		"access-token":  token.AccessToken,
		"refresh-token": token.RefreshToken,
	}

	ctx.JSON(http.StatusOK, jsonMsg)
}
