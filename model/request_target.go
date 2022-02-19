//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestTargetRepository interface {
	GetRequestTargets(ctx context.Context, requestID uuid.UUID) ([]*RequestTarget, error)
}

type RequestTarget struct {
	ID        uuid.UUID
	Target    string
	PaidAt    *time.Time
	CreatedAt time.Time
}
