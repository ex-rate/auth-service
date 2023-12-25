package auth

import (
	"context"

	"github.com/ex-rate/auth-service/internal/entities"
	token "github.com/ex-rate/auth-service/internal/service/token"
	"github.com/google/uuid"

	schema "github.com/ex-rate/auth-service/internal/schemas"
)

type AuthService struct {
	authRepo authRepo
	token    *token.Token
}

//go:generate mockgen -source auth.go -destination ../../mocks/auth_repo.go -package mock_service
type authRepo interface {
	GetUserID(ctx context.Context, username string) (uuid.UUID, error)
	GetHashPassword(ctx context.Context, auth schema.AuthWithPassword) (string, error)
}

func New(authRepo authRepo, token *token.Token) *AuthService {
	return &AuthService{authRepo: authRepo, token: token}
}

// WithCode проводит авторизацию пользователя по коду: проверяет введенные данные и возвращает токены и ошибку
func (s *AuthService) WithCode(ctx context.Context, auth schema.AuthWithCode) (*schema.Token, error) {
	// проверяем существование пользователя в базе
	userID, err := s.authRepo.GetUserID(ctx, auth.Username)
	if err != nil {
		return nil, err
	}

	// проверяем код
	err = s.checkCode(auth)
	if err != nil {
		return nil, err
	}

	user := entities.Token{
		UserID:   userID,
		Username: auth.Username,
	}

	// выдаем токен
	return s.token.GenerateTokens(ctx, user)
}

// WithCode проводит авторизацию пользователя по паролю: проверяет введенные данные и возвращает токены и ошибку
func (s *AuthService) WithPassword(ctx context.Context, user schema.AuthWithPassword) (*schema.Token, error) {
	// проверяем существование пользователя в базе
	userID, err := s.authRepo.GetUserID(ctx, user.Username)
	if err != nil {
		return nil, err
	}

	user.UserID = userID

	// проверяем пароль
	err = s.checkPassword(ctx, user)
	if err != nil {
		return nil, err
	}

	userForToken := entities.Token{
		UserID:   userID,
		Username: user.Username,
	}

	// выдаем токен
	return s.token.GenerateTokens(ctx, userForToken)
}
