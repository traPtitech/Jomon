//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestTargetRepository interface {
	GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*TargetDetail, error)
}

type Target struct {
	Target string
	Amount int
}
type TargetDetail struct {
	ID        uuid.UUID
	Target    string
	Amount    int
	PaidAt    *time.Time
	CreatedAt time.Time
}
