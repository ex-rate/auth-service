package postgres

import (
	"github.com/jmoiron/sqlx"
)

func Connect(dbString string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dbString)
	if err != nil {
		return nil, err
	}

	return db, nil
}
