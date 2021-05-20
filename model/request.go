package model

import (
	"time"

	"github.com/google/uuid"
)

type RequestRepository interface {
}

type Request struct {
	ID        uuid.UUID
	Amount    int
	CreatedAt time.Time
}
