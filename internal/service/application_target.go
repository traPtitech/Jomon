//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/internal/nulltime"
)

type ApplicationTargetRepository interface {
	GetApplicationTargets(
		ctx context.Context,
		applicationID uuid.UUID,
	) ([]*ApplicationTargetDetail, error)
}

type ApplicationTargetDetail struct {
	ID        uuid.UUID
	Target    uuid.UUID
	Amount    int
	PaidAt    nulltime.NullTime
	CreatedAt time.Time
}

type ApplicationTarget struct {
	Target uuid.UUID
	Amount int
}
