//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestRepository interface {
	GetRequests(ctx context.Context, query RequestQuery) ([]*RequestResponse, error)
	CreateRequest(
		ctx context.Context, title string, content string,
		tags []*Tag, targets []*RequestTarget, group *Group, userID uuid.UUID,
	) (*RequestDetail, error)
	GetRequest(ctx context.Context, requestID uuid.UUID) (*RequestDetail, error)
	UpdateRequest(
		ctx context.Context, requestID uuid.UUID, title string, content string,
		tags []*Tag, targets []*RequestTarget, group *Group,
	) (*RequestDetail, error)
}

type Request struct {
	ID        uuid.UUID
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestResponse struct {
	ID        uuid.UUID
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
	Title     string
	Content   string
	Tags      []*Tag
	Targets   []*RequestTargetDetail
	Statuses  []*RequestStatus
	Group     *Group
}

type RequestDetail struct {
	ID        uuid.UUID
	Status    Status
	Title     string
	Content   string
	Comments  []*Comment
	Files     []*uuid.UUID
	Statuses  []*RequestStatus
	Tags      []*Tag
	Targets   []*RequestTargetDetail
	Group     *Group
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
}

type RequestQuery struct {
	Sort      *string
	Target    *uuid.UUID
	Status    *string
	Since     *time.Time
	Until     *time.Time
	Limit     int
	Offset    int
	Tag       *string
	Group     *string
	CreatedBy *uuid.UUID
}
