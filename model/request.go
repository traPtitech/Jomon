package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RequestRepository interface {
	GetRequests(ctx context.Context, query RequestQuery) ([]*RequestResponse, error)
	CreateRequest(ctx context.Context, amount int, title string, content string, tags []*Tag, group *Group, files []*File) (*RequestDetail, error)
	/*
		GetRequest(ctx context.Context, requestID uuid.UUID) (RequestDetail, error)
		UpdateRequest(ctx context.Context, requestID uuid.UUID, amount int, title string, content string, tags []*Tag, group Group) (RequestDetail, error)
	*/
}

type Request struct {
	ID        uuid.UUID
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestResponse struct {
	ID        uuid.UUID
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
	Amount    int
	Title     string
	Content   string
	Tags      []*Tag
	Group     Group
}

type RequestDetail struct {
	ID        uuid.UUID
	Amount    int
	Title     string
	Content   string
	Comments  []*Comment
	Files     []*uuid.UUID
	Statuses  []*RequestStatus
	Tags      []*Tag
	Group     Group
	CreatedAt time.Time
	UpdatedAt time.Time
	CreatedBy uuid.UUID
}

type RequestQuery struct {
	Sort   *string
	Target *string
	Year   *int
	Since  *time.Time
	Until  *time.Time
	Tag    *string
	Group  *string
}
