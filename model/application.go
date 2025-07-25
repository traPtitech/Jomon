//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ApplicationRepository interface {
	GetApplications(ctx context.Context, query ApplicationQuery) ([]*ApplicationResponse, error)
	CreateApplication(
		ctx context.Context, title string, content string,
		tags []*Tag, targets []*ApplicationTarget, group *Group, userID uuid.UUID,
	) (*ApplicationDetail, error)
	GetApplication(ctx context.Context, applicationID uuid.UUID) (*ApplicationDetail, error)
	UpdateApplication(
		ctx context.Context, applicationID uuid.UUID, title string, content string,
		tags []*Tag, targets []*ApplicationTarget, group *Group,
	) (*ApplicationDetail, error)
}

type Application struct {
	ID        uuid.UUID
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ApplicationResponse struct {
	ID        uuid.UUID
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
	Title     string
	Content   string
	Tags      []*Tag
	Targets   []*ApplicationTargetDetail
	Statuses  []*ApplicationStatus
	Group     *Group
}

type ApplicationDetail struct {
	ID        uuid.UUID
	Status    Status
	Title     string
	Content   string
	Comments  []*Comment
	Files     []uuid.UUID
	Statuses  []*ApplicationStatus
	Tags      []*Tag
	Targets   []*ApplicationTargetDetail
	Group     *Group
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
}

type ApplicationQuery struct {
	Sort      *string
	Target    uuid.UUID
	Status    *string
	Since     time.Time
	Until     time.Time
	Limit     int
	Offset    int
	Tag       *string
	Group     *string
	CreatedBy uuid.UUID
}
