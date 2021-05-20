package model

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID
}

type UserRepository interface {
}
