package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/comment"
	"github.com/traPtitech/Jomon/ent/request"
)

func (repo *EntRepository) GetComments(ctx context.Context, requestID uuid.UUID) ([]*Comment, error) {
	comments, err := repo.client.Comment.
		Query().
		Where(
			comment.HasRequestWith(
				request.ID(requestID),
			),
		).
		All(ctx)
	if err != nil {
		return nil, err
	}
	modelcomments := []*Comment{}
	for _, comment := range comments {
		modelcomments = append(modelcomments, ConvertEntCommentToModelComment(comment))
	}
	return modelcomments, nil
}

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

// TODO: add edge to request
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

func ConvertEntCommentToModelComment(comment *ent.Comment) *Comment {
	return &Comment{
		ID:        comment.ID,
		User:      comment.Edges.User.ID,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
