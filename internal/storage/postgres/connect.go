package postgres

import (
	"github.com/jmoiron/sqlx"
)

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=foo dbname=bar sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}
