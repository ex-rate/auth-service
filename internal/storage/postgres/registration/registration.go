package registration

import (
	"context"

	api_errors "github.com/ex-rate/auth-service/internal/errors"
	schema "github.com/ex-rate/auth-service/internal/schemas"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type registrationRepo struct {
	conn *sqlx.DB
}

func New(conn *sqlx.DB) *registrationRepo {
	return &registrationRepo{conn: conn}
}

func (db *registrationRepo) CreateUser(ctx context.Context, user schema.Registration) error {
	var userID uuid.UUID

	tx := db.conn.MustBegin()
	rows, err := tx.QueryxContext(ctx, `insert into auth.users (username, hash_password, fullname) values($1, $2, $3) returning user_id`, user.Username, user.HashedPassword, user.FullName)
	if err != nil {
		switch t := err.(type) {
		case *pq.Error: // refactor
			if t.Code == "23505" {
				if t.Constraint == "users_username_key" {
					return api_errors.ErrUsernameAlreadyExists
				}
			}
		default:
			return err
		}
		return err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return err
		}

	}

	err = db.processPhoneOrEmail(ctx, tx, userID, user)
	if err != nil {
		return handleDbError(err)
	}

	return tx.Commit()
}

func handleDbError(err error) error {
	switch t := err.(type) {
	case *pq.Error: // refactor
		if t.Code == "23505" {
			switch t.Constraint {
			case "emails_email_key":
				return api_errors.ErrEmailAlreadyExists
			case "phone_numbers_phone_number_key":
				return api_errors.ErrPhoneAlreadyExists
			}
		}
	default:
		return err
	}

	return err
}

func (db *registrationRepo) processPhoneOrEmail(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, user schema.Registration) error {
	if user.Email != "" {
		return db.insertEmail(ctx, tx, userID, user)
	}

	return db.insertPhone(ctx, tx, userID, user)
}

func (db *registrationRepo) insertEmail(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, user schema.Registration) error {
	_, err := tx.ExecContext(ctx, `insert into auth.emails (user_id, email) values($1, $2)`, userID, user.Email)
	return err
}

func (db *registrationRepo) insertPhone(ctx context.Context, tx *sqlx.Tx, userID uuid.UUID, user schema.Registration) error {
	_, err := tx.ExecContext(ctx, `insert into auth.phone_numbers (user_id, phone_number) values($1, $2)`, userID, user.PhoneNumber)
	return err
}
