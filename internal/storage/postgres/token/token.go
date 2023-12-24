package token

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ex-rate/auth-service/internal/entities"
	api_errors "github.com/ex-rate/auth-service/internal/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type tokenRepo struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *tokenRepo {
	return &tokenRepo{conn: conn}
}

func (db *tokenRepo) CreateToken(ctx context.Context, user *entities.Token) error {
	q := `insert into auth.refresh_tokens (user_id, token, expiration_time) 
	values(:user_id, :token, :expiration_time)
	on conflict (user_id) do update 
	set 
	token = :token,
	expiration_time = :expiration_time`

	tx := db.conn.MustBegin()

	_, err := tx.NamedExecContext(ctx, q, user)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *tokenRepo) CheckToken(ctx context.Context, token *entities.Token) error {
	var id int

	fmt.Println(token.RefreshToken)

	err := db.conn.GetContext(ctx, &id, `select id from auth.refresh_tokens where token = $1`, token.RefreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return api_errors.ErrTokenNotExists
		}
		return err
	}

	log.Println("check token ID: ", id)

	return nil
}

func (db *tokenRepo) GetUserID(ctx context.Context, username string) (uuid.UUID, error) {
	var id uuid.UUID

	err := db.conn.GetContext(ctx, &id, `select user_id from auth.users where username = $1`, username)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}
