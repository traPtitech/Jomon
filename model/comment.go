//go:generate go tool mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
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
	GetComments(ctx context.Context, applicationID uuid.UUID) ([]*Comment, error)
	CreateComment(
		ctx context.Context, comment string, applicationID uuid.UUID, userID uuid.UUID,
	) (*Comment, error)
	UpdateComment(
		ctx context.Context, comment string, applicationID uuid.UUID, commentID uuid.UUID,
	) (*Comment, error)
	DeleteComment(ctx context.Context, applicationID uuid.UUID, commentID uuid.UUID) error
}
