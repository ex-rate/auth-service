package entities

import (
	"github.com/google/uuid"
)

type Token struct {
	UserID       uuid.UUID `db:"user_id"`
	Username     string    `db:"-"`
	RefreshToken string    `db:"token"`
	ExpTime      float64   `db:"expiration_time"`
}

type RestoreToken struct {
	RefreshToken string `json:"refresh-token"`
	AccessToken  string
}
