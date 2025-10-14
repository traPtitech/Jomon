package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/ent/application"
	"github.com/traPtitech/Jomon/internal/ent/applicationstatus"
	"github.com/traPtitech/Jomon/internal/ent/applicationtarget"
	"github.com/traPtitech/Jomon/internal/ent/tag"
	"github.com/traPtitech/Jomon/internal/ent/user"
)

func (repo *EntRepository) GetApplications(
	ctx context.Context, query ApplicationQuery,
) ([]*ApplicationResponse, error) {
	// Querying
	var applicationsq *ent.ApplicationQuery
	var err error
	if query.Sort == nil || *query.Sort == "" || *query.Sort == "created_at" {
		applicationsq = repo.client.Application.
			Query().
			WithTag().
			WithStatus(func(q *ent.ApplicationStatusQuery) {
				q.Order(ent.Desc(applicationstatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Desc(application.FieldCreatedAt))
	} else if *query.Sort == "-created_at" {
		applicationsq = repo.client.Application.
			Query().
			WithTag().
			WithStatus(func(q *ent.ApplicationStatusQuery) {
				q.Order(ent.Desc(applicationstatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Asc(application.FieldCreatedAt))
	} else if *query.Sort == "title" {
		applicationsq = repo.client.Application.
			Query().
			WithTag().
			WithStatus(func(q *ent.ApplicationStatusQuery) {
				q.Order(ent.Desc(applicationstatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Asc(application.FieldTitle))
	} else if *query.Sort == "-title" {
		applicationsq = repo.client.Application.
			Query().
			WithTag().
			WithStatus(func(q *ent.ApplicationStatusQuery) {
				q.Order(ent.Desc(applicationstatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Desc(application.FieldTitle))
	}

	if query.Target != uuid.Nil {
		applicationsq = applicationsq.
			Where(
				application.HasTargetWith(
					applicationtarget.HasUserWith(
						user.IDEQ(query.Target),
					),
				),
			)
	}

	if query.Status != nil && *query.Status != "" {
		applicationsq = applicationsq.
			Where(
				application.HasStatusWith(
					applicationstatus.StatusEQ(applicationstatus.Status(*query.Status)),
				),
			)
	}

	if !query.Since.IsZero() {
		applicationsq = applicationsq.
			Where(application.CreatedAtGTE(query.Since))
	}

	if !query.Until.IsZero() {
		applicationsq = applicationsq.
			Where(application.CreatedAtLT(query.Until))
	}

	if query.Tag != nil && *query.Tag != "" {
		applicationsq = applicationsq.
			Where(
				application.HasTagWith(
					tag.NameEQ(*query.Tag),
				),
			)
	}

	applicationsq = applicationsq.Limit(query.Limit).Offset(query.Offset)

	if query.CreatedBy != uuid.Nil {
		applicationsq = applicationsq.
			Where(
				application.HasUserWith(
					user.IDEQ(query.CreatedBy),
				),
			)
	}

	applications, err := applicationsq.All(ctx)
	if err != nil {
		return nil, err
	}

	reqres := lo.Map(applications, func(r *ent.Application, _ int) *ApplicationResponse {
		return convertEntApplicationResponseToModelApplicationResponse(
			r, r.Edges.Tag, r.Edges.Status[0], r.Edges.User)
	})

	return reqres, nil
}

func (repo *EntRepository) CreateApplication(
	ctx context.Context, title string, content string,
	tags []*Tag, targets []*ApplicationTarget, userID uuid.UUID,
) (*ApplicationDetail, error) {
	tx, err := repo.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()
	tagIDs := lo.Map(tags, func(t *Tag, _ int) uuid.UUID {
		return t.ID
	})
	created, err := tx.Client().Application.
		Create().
		SetTitle(title).
		SetContent(content).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetUserID(userID).
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	t, err := created.QueryUser().Select(user.FieldID).First(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	s, err := tx.Client().ApplicationStatus.
		Create().
		SetStatus(applicationstatus.StatusSubmitted).
		SetCreatedAt(time.Now()).
		SetApplication(created).
		SetUser(t).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	status, err := tx.Client().ApplicationStatus.
		Query().
		Where(applicationstatus.IDEQ(s.ID)).
		WithUser().
		First(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	ts, err := repo.createApplicationTargets(ctx, tx, created.ID, targets)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	statuses := []*ApplicationStatus{convertEntApplicationStatusToModelApplicationStatus(status)}
	reqdetail := &ApplicationDetail{
		ID:        created.ID,
		Status:    convertEntApplicationStatusToModelStatus(&status.Status),
		Title:     created.Title,
		Content:   created.Content,
		Tags:      tags,
		Targets:   ts,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
		CreatedBy: t.ID,
		Comments:  []*Comment{},
		Statuses:  statuses,
		Files:     []uuid.UUID{},
	}
	return reqdetail, nil
}

func (repo *EntRepository) GetApplication(
	ctx context.Context, applicationID uuid.UUID,
) (*ApplicationDetail, error) {
	r, err := repo.client.Application.
		Query().
		Where(application.IDEQ(applicationID)).
		WithTag().
		WithTarget(func(q *ent.ApplicationTargetQuery) {
			q.WithUser()
		}).
		WithStatus(func(q *ent.ApplicationStatusQuery) {
			q.Order(ent.Desc(applicationstatus.FieldCreatedAt))
			q.WithUser()
		}).
		WithUser().
		WithComment().
		WithFile().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	tags := lo.Map(r.Edges.Tag, func(t *ent.Tag, _ int) *Tag {
		return ConvertEntTagToModelTag(t)
	})
	targets := lo.Map(
		r.Edges.Target,
		func(target *ent.ApplicationTarget, _ int) *ApplicationTargetDetail {
			return ConvertEntApplicationTargetToModelApplicationTargetDetail(target)
		},
	)
	comments := lo.Map(r.Edges.Comment, func(c *ent.Comment, _ int) *Comment {
		return ConvertEntCommentToModelComment(c, c.Edges.User.ID)
	})
	statuses := lo.Map(
		r.Edges.Status,
		func(status *ent.ApplicationStatus, _ int) *ApplicationStatus {
			return convertEntApplicationStatusToModelApplicationStatus(status)
		},
	)
	files := lo.Map(r.Edges.File, func(f *ent.File, _ int) uuid.UUID {
		return f.ID
	})
	reqdetail := &ApplicationDetail{
		ID:        r.ID,
		Status:    convertEntApplicationStatusToModelStatus(&r.Edges.Status[0].Status),
		Title:     r.Title,
		Content:   r.Content,
		Tags:      tags,
		Targets:   targets,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		CreatedBy: r.Edges.User.ID,
		Comments:  comments,
		Statuses:  statuses,
		Files:     files,
	}
	return reqdetail, nil
}

func (repo *EntRepository) UpdateApplication(
	ctx context.Context, applicationID uuid.UUID, title string, content string,
	tags []*Tag, targets []*ApplicationTarget,
) (*ApplicationDetail, error) {
	tx, err := repo.client.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()
	tagIDs := lo.Map(tags, func(t *Tag, _ int) uuid.UUID {
		return t.ID
	})
	updated, err := tx.Client().Application.
		UpdateOneID(applicationID).
		SetTitle(title).
		SetContent(content).
		ClearTag().
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	_, err = tx.Client().Application.
		UpdateOneID(applicationID).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	entstatuses, err := updated.QueryStatus().
		WithUser().
		Order(ent.Desc(applicationstatus.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	status := entstatuses[0]
	u, err := updated.QueryUser().Select(user.FieldID).First(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	enttags, err := updated.QueryTag().All(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	modeltags := lo.Map(enttags, func(enttag *ent.Tag, _ int) *Tag {
		return ConvertEntTagToModelTag(enttag)
	})

	err = repo.deleteApplicationTargets(ctx, tx, applicationID)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	modeltargets, err := repo.createApplicationTargets(ctx, tx, applicationID, targets)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	entcomments, err := updated.QueryComment().
		WithUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	comments := lo.Map(entcomments, func(c *ent.Comment, _ int) *Comment {
		return ConvertEntCommentToModelComment(c, c.Edges.User.ID)
	})
	statuses := lo.Map(entstatuses, func(s *ent.ApplicationStatus, _ int) *ApplicationStatus {
		return convertEntApplicationStatusToModelApplicationStatus(s)
	})
	entfiles, err := updated.QueryFile().All(ctx)
	if err != nil {
		return nil, err
	}
	files := lo.Map(entfiles, func(f *ent.File, _ int) uuid.UUID {
		return f.ID
	})
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	reqdetail := &ApplicationDetail{
		ID:        updated.ID,
		Status:    convertEntApplicationStatusToModelStatus(&status.Status),
		Title:     updated.Title,
		Content:   updated.Content,
		Tags:      modeltags,
		Targets:   modeltargets,
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
		CreatedBy: u.ID,
		Comments:  comments,
		Statuses:  statuses,
		Files:     files,
	}
	return reqdetail, nil
}

func convertEntApplicationResponseToModelApplicationResponse(
	application *ent.Application, tags []*ent.Tag,
	status *ent.ApplicationStatus, user *ent.User,
) *ApplicationResponse {
	if application == nil {
		return nil
	}
	modeltags := lo.Map(tags, func(t *ent.Tag, _ int) *Tag {
		return ConvertEntTagToModelTag(t)
	})
	return &ApplicationResponse{
		ID:        application.ID,
		Status:    convertEntApplicationStatusToModelStatus(&status.Status),
		CreatedAt: application.CreatedAt,
		UpdatedAt: application.UpdatedAt,
		CreatedBy: user.ID,
		Title:     application.Title,
		Content:   application.Content,
		Tags:      modeltags,
	}
}
