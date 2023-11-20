package entities

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID `db:"user_id"`
	Username       string    `db:"username"`
	HashedPassword string    `db:"hash_password"`
	Email          string    `db:"email"`
	PhoneNumber    string    `db:"phone_number"`
	FullName       string    `db:"fullname"`
}
