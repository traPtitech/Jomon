package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
)

func (repo *EntRepository) CreateComment(ctx context.Context, comment string, requestID uuid.UUID, userID uuid.UUID) (*Comment, error) {
	created, err := repo.client.Comment.
		Create().
		SetComment(comment).
		SetRequestID(requestID).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntCommentToModelComment(created), nil
}

func (repo *EntRepository) UpdateComment(ctx context.Context, comment string, requestID uuid.UUID, commentID uuid.UUID) (*Comment, error) {
	updated, err := repo.client.Comment.
		UpdateOneID(commentID).
		SetComment(comment).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntCommentToModelComment(updated), nil
}

func ConvertEntCommentToModelComment(entcomment *ent.Comment) *Comment {
	return &Comment{
		ID:        entcomment.ID,
		User:      entcomment.Edges.User.ID,
		Comment:   entcomment.Comment,
		CreatedAt: entcomment.CreatedAt,
		UpdatedAt: entcomment.UpdatedAt,
	}
}
