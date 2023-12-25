package schema

import "github.com/google/uuid"

type AuthWithCode struct {
	Username string `json:"username"`
	Code     int    `json:"code"`
}

type AuthWithPassword struct {
	Username string `json:"username"`
	UserID   uuid.UUID
	Password string `json:"password"`
}
