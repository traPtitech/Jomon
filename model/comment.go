//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID        uuid.UUID
	User      uuid.UUID
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentRepository interface {
	CreateComment(ctx context.Context, comment string, requestID uuid.UUID, userID uuid.UUID) (*Comment, error)
	UpdateComment(ctx context.Context, comment string, requestID uuid.UUID, commentID uuid.UUID) (*Comment, error)
	DeleteComment(ctx context.Context, requestID uuid.UUID, commentID uuid.UUID) error
}
