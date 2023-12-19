package service

import (
	"context"
	"time"

	"github.com/ex-rate/auth-service/internal/entities"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/golang-jwt/jwt"
)

// GenerateToken генерирует новые токены: refresh и access
func (s *token) GenerateToken(ctx context.Context, user entities.Token) (*schema.Token, error) {
	accessToken, err := s.accessToken(user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.refreshToken(user.Username)
	if err != nil {
		return nil, err
	}

	userToken := &schema.Token{
		RefreshToken: refreshToken.RefreshToken,
		AccessToken:  accessToken,
	}

	user.RefreshToken = refreshToken.RefreshToken
	user.ExpTime = refreshToken.ExpTime

	err = s.tokenRepo.CreateToken(ctx, &user)
	if err != nil {
		return nil, err
	}

	return userToken, nil
}

// accessToken создает новый access token
func (s *token) accessToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	key := []byte(s.secretKey)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix() // поменять потом
	claims["authorized"] = true
	claims["user"] = username

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// refreshToken создает новый refresh token
func (s *token) refreshToken(username string) (*entities.Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	key := []byte(s.secretKey)

	day := time.Hour * 24
	exprTime := time.Now().Add(30 * day).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exprTime
	claims["authorized"] = true
	claims["user"] = username

	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}

	entity := &entities.Token{
		RefreshToken: tokenString,
		ExpTime:      float64(exprTime),
	}

	return entity, nil
}
