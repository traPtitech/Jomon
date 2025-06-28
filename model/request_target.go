//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestTargetRepository interface {
	GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*RequestTargetDetail, error)
}

type RequestTargetDetail struct {
	ID        uuid.UUID
	Target    uuid.UUID
	Amount    int
	PaidAt    *time.Time
	CreatedAt time.Time
}

type RequestTarget struct {
	Target uuid.UUID
	Amount int
}
