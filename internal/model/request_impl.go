package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requeststatus"
	"github.com/traPtitech/Jomon/ent/requesttarget"
	"github.com/traPtitech/Jomon/ent/tag"
	"github.com/traPtitech/Jomon/ent/user"
)

func (repo *EntRepository) GetRequests(
	ctx context.Context, query RequestQuery,
) ([]*RequestResponse, error) {
	// Querying
	var requestsq *ent.RequestQuery
	var err error
	if query.Sort == nil || *query.Sort == "" || *query.Sort == "created_at" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Desc(request.FieldCreatedAt))
	} else if *query.Sort == "-created_at" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Asc(request.FieldCreatedAt))
	} else if *query.Sort == "title" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Asc(request.FieldTitle))
	} else if *query.Sort == "-title" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Desc(request.FieldTitle))
	}

	if query.Target != uuid.Nil {
		requestsq = requestsq.
			Where(
				request.HasTargetWith(
					requesttarget.HasUserWith(
						user.IDEQ(query.Target),
					),
				),
			)
	}

	if query.Status != nil && *query.Status != "" {
		requestsq = requestsq.
			Where(
				request.HasStatusWith(
					requeststatus.StatusEQ(requeststatus.Status(*query.Status)),
				),
			)
	}

	if !query.Since.IsZero() {
		requestsq = requestsq.
			Where(request.CreatedAtGTE(query.Since))
	}

	if !query.Until.IsZero() {
		requestsq = requestsq.
			Where(request.CreatedAtLT(query.Until))
	}

	if query.Tag != nil && *query.Tag != "" {
		requestsq = requestsq.
			Where(
				request.HasTagWith(
					tag.NameEQ(*query.Tag),
				),
			)
	}

	requestsq = requestsq.Limit(query.Limit).Offset(query.Offset)

	if query.CreatedBy != uuid.Nil {
		requestsq = requestsq.
			Where(
				request.HasUserWith(
					user.IDEQ(query.CreatedBy),
				),
			)
	}

	requests, err := requestsq.All(ctx)
	if err != nil {
		return nil, err
	}

	reqres := lo.Map(requests, func(r *ent.Request, _ int) *RequestResponse {
		return convertEntRequestResponseToModelRequestResponse(
			r, r.Edges.Tag, r.Edges.Status[0], r.Edges.User)
	})

	return reqres, nil
}

func (repo *EntRepository) CreateRequest(
	ctx context.Context, title string, content string,
	tags []*Tag, targets []*RequestTarget, userID uuid.UUID,
) (*RequestDetail, error) {
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
	created, err := tx.Client().Request.
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
	s, err := tx.Client().RequestStatus.
		Create().
		SetStatus(requeststatus.StatusSubmitted).
		SetCreatedAt(time.Now()).
		SetRequest(created).
		SetUser(t).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	status, err := tx.Client().RequestStatus.
		Query().
		Where(requeststatus.IDEQ(s.ID)).
		WithUser().
		First(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	ts, err := repo.createRequestTargets(ctx, tx, created.ID, targets)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	statuses := []*RequestStatus{convertEntRequestStatusToModelRequestStatus(status)}
	reqdetail := &RequestDetail{
		ID:        created.ID,
		Status:    convertEntRequestStatusToModelStatus(&status.Status),
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

func (repo *EntRepository) GetRequest(
	ctx context.Context, requestID uuid.UUID,
) (*RequestDetail, error) {
	r, err := repo.client.Request.
		Query().
		Where(request.IDEQ(requestID)).
		WithTag().
		WithTarget(func(q *ent.RequestTargetQuery) {
			q.WithUser()
		}).
		WithStatus(func(q *ent.RequestStatusQuery) {
			q.Order(ent.Desc(requeststatus.FieldCreatedAt))
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
		func(target *ent.RequestTarget, _ int) *RequestTargetDetail {
			return ConvertEntRequestTargetToModelRequestTargetDetail(target)
		},
	)
	comments := lo.Map(r.Edges.Comment, func(c *ent.Comment, _ int) *Comment {
		return ConvertEntCommentToModelComment(c, c.Edges.User.ID)
	})
	statuses := lo.Map(r.Edges.Status, func(status *ent.RequestStatus, _ int) *RequestStatus {
		return convertEntRequestStatusToModelRequestStatus(status)
	})
	files := lo.Map(r.Edges.File, func(f *ent.File, _ int) uuid.UUID {
		return f.ID
	})
	reqdetail := &RequestDetail{
		ID:        r.ID,
		Status:    convertEntRequestStatusToModelStatus(&r.Edges.Status[0].Status),
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

func (repo *EntRepository) UpdateRequest(
	ctx context.Context, requestID uuid.UUID, title string, content string,
	tags []*Tag, targets []*RequestTarget,
) (*RequestDetail, error) {
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
	updated, err := tx.Client().Request.
		UpdateOneID(requestID).
		SetTitle(title).
		SetContent(content).
		ClearTag().
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	_, err = tx.Client().Request.
		UpdateOneID(requestID).
		Save(ctx)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}

	entstatuses, err := updated.QueryStatus().
		WithUser().
		Order(ent.Desc(requeststatus.FieldCreatedAt)).
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

	err = repo.deleteRequestTargets(ctx, tx, requestID)
	if err != nil {
		err = RollbackWithError(tx, err)
		return nil, err
	}
	modeltargets, err := repo.createRequestTargets(ctx, tx, requestID, targets)
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
	statuses := lo.Map(entstatuses, func(s *ent.RequestStatus, _ int) *RequestStatus {
		return convertEntRequestStatusToModelRequestStatus(s)
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

	reqdetail := &RequestDetail{
		ID:        updated.ID,
		Status:    convertEntRequestStatusToModelStatus(&status.Status),
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

func convertEntRequestResponseToModelRequestResponse(
	request *ent.Request, tags []*ent.Tag,
	status *ent.RequestStatus, user *ent.User,
) *RequestResponse {
	if request == nil {
		return nil
	}
	modeltags := lo.Map(tags, func(t *ent.Tag, _ int) *Tag {
		return ConvertEntTagToModelTag(t)
	})
	return &RequestResponse{
		ID:        request.ID,
		Status:    convertEntRequestStatusToModelStatus(&status.Status),
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
		CreatedBy: user.ID,
		Title:     request.Title,
		Content:   request.Content,
		Tags:      modeltags,
	}
}
