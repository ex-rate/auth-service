package service

import (
	"time"

	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/golang-jwt/jwt"
)

func (s *Service) GenerateToken(user schema.Registration) (string, error) {
	key := []byte(s.secretKey)

	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute) // поменять потом
	claims["authorized"] = true
	claims["user"] = user.Username

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
