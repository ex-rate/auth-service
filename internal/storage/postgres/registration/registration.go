package registration

import (
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/jmoiron/sqlx"
)

type RegistrationRepo struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *RegistrationRepo {
	return &RegistrationRepo{conn: conn}
}

func (db *RegistrationRepo) CreateUser(reg schema.Registration) error {
	tx := db.conn.MustBegin()
	_, err := tx.NamedExec(`insert into auth.users (:username, :hash_password, :email, :phone_number, :fullname) values($1, $2, $3, $4, $5)`, &reg)
	if err != nil {
		return err
	}

	return tx.Commit()
}
