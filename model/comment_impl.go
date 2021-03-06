package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/comment"
	"github.com/traPtitech/Jomon/ent/request"
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

func (repo *EntRepository) DeleteComment(ctx context.Context, requestID uuid.UUID, commentID uuid.UUID) error {
	comment, err := repo.client.Comment.
		Query().
		Where(
			comment.HasRequestWith(
				request.ID(requestID),
			),
		).
		Where(comment.IDEQ(commentID)).
		Only(ctx)
	if err != nil {
		return err
	}
	err = repo.client.Comment.
		DeleteOne(comment).
		Exec(ctx)
	return err
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
