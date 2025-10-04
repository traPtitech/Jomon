package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/application"
	"github.com/traPtitech/Jomon/ent/comment"
)

func (repo *EntRepository) GetComments(
	ctx context.Context, applicationID uuid.UUID,
) ([]*Comment, error) {
	_, err := repo.client.Application.
		Query().
		Where(application.IDEQ(applicationID)).
		First(ctx)
	if err != nil {
		return nil, err
	}

	comments, err := repo.client.Comment.
		Query().
		Where(
			comment.HasApplicationWith(
				application.ID(applicationID),
			),
		).
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modelcomments := lo.Map(comments, func(c *ent.Comment, _ int) *Comment {
		return ConvertEntCommentToModelComment(c, c.Edges.User.ID)
	})
	return modelcomments, nil
}

func (repo *EntRepository) CreateComment(
	ctx context.Context, comment string, applicationID uuid.UUID, userID uuid.UUID,
) (*Comment, error) {
	created, err := repo.client.Comment.
		Create().
		SetComment(comment).
		SetApplicationID(applicationID).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntCommentToModelComment(created, userID), nil
}

func (repo *EntRepository) UpdateComment(
	ctx context.Context, commentContent string, applicationID uuid.UUID, commentID uuid.UUID,
) (*Comment, error) {
	updated, err := repo.client.Comment.
		UpdateOneID(commentID).
		SetComment(commentContent).
		SetUpdatedAt(time.Now()).
		SetApplicationID(applicationID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	updatedWithUser, err := repo.client.Comment.
		Query().
		Where(comment.IDEQ(commentID)).
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntCommentToModelComment(updated, updatedWithUser.Edges.User.ID), nil
}

func (repo *EntRepository) DeleteComment(
	ctx context.Context, applicationID uuid.UUID, commentID uuid.UUID,
) error {
	c, err := repo.client.Comment.
		Query().
		Where(
			comment.HasApplicationWith(
				application.ID(applicationID),
			),
		).
		Where(comment.IDEQ(commentID)).
		Only(ctx)
	if err != nil {
		return err
	}
	err = repo.client.Comment.
		DeleteOne(c).
		Exec(ctx)
	return err
}

func ConvertEntCommentToModelComment(comment *ent.Comment, userID uuid.UUID) *Comment {
	return &Comment{
		ID:        comment.ID,
		User:      userID,
		Comment:   comment.Comment,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: comment.UpdatedAt,
	}
}
