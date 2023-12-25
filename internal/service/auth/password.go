package auth

import (
	"context"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	token "github.com/ex-rate/auth-service/internal/service/token"
)

// checkPassword проводит проверку пароля
func (s *AuthService) checkPassword(ctx context.Context, user schema.AuthWithPassword) error {
	hash, err := s.authRepo.GetHashPassword(ctx, user)
	if err != nil {
		return err
	}

	if !token.CheckPasswordHash(user.Password, hash) {
		return api_errors.ErrIncorrectPassword
	}

	return nil
}
