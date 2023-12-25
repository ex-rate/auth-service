package auth

import (
	"context"
	"database/sql"
	"errors"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type authRepo struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *authRepo {
	return &authRepo{conn: conn}
}

func (db *authRepo) GetUserID(ctx context.Context, username string) (uuid.UUID, error) {
	var id uuid.UUID

	err := db.conn.GetContext(ctx, &id, `select user_id from auth.users where username = $1`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, api_errors.ErrUserNotExists
		}
		return uuid.UUID{}, err
	}

	return id, nil
}

func (db *authRepo) GetHashPassword(ctx context.Context, user schema.AuthWithPassword) (string, error) {
	var hash string

	err := db.conn.GetContext(ctx, &hash, `select hash_password from auth.users where username = $1 and user_id = $2`, user.Username, user.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", api_errors.ErrUserNotExists
		}
		return "", err
	}

	return hash, nil
}
