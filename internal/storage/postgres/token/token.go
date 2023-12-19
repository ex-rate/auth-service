package token

import (
	"context"
	"fmt"
	"log"

	"github.com/ex-rate/auth-service/internal/entities"
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
		return err
	}

	log.Println("check token ID: ", id)

	return nil
}
