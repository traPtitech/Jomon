package model

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/group"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requeststatus"
	"github.com/traPtitech/Jomon/ent/requesttarget"
	"github.com/traPtitech/Jomon/ent/tag"
	"github.com/traPtitech/Jomon/ent/user"
)

func (repo *EntRepository) GetRequests(ctx context.Context, query RequestQuery) ([]*RequestResponse, error) {
	// Querying
	var requestsq *ent.RequestQuery
	var err error
	if query.Sort == nil || *query.Sort == "" || *query.Sort == "created_at" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Desc(request.FieldCreatedAt))
	} else if *query.Sort == "-created_at" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Asc(request.FieldCreatedAt))
	} else if *query.Sort == "title" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Desc(request.FieldTitle))
	} else if *query.Sort == "-title" {
		requestsq = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt))
			}).
			WithUser().
			Order(ent.Asc(request.FieldTitle))
	}

	if query.Target != nil && *query.Target != "" {
		requestsq = requestsq.
			Where(
				request.HasTargetWith(
					requesttarget.TargetEQ(*query.Target),
				),
			)
	}

	if query.Year != nil && *query.Year != 0 {
		requestsq = requestsq.
			Where(request.CreatedAtGTE(time.Date(*query.Year, 4, 1, 0, 0, 0, 0, time.Local))).
			Where(request.CreatedAtLT(time.Date(*query.Year+1, 4, 1, 0, 0, 0, 0, time.Local)))
	}

	if query.Since != nil && !(*query.Since).IsZero() {
		requestsq = requestsq.
			Where(request.CreatedAtGTE(*query.Since))
	}

	if query.Until != nil && !(*query.Until).IsZero() {
		requestsq = requestsq.
			Where(request.CreatedAtLT(*query.Until))
	}

	if query.Tag != nil && *query.Tag != "" {
		requestsq = requestsq.
			Where(
				request.HasTagWith(
					tag.NameEQ(*query.Tag),
				),
			)
	}

	if query.Group != nil && *query.Group != "" {
		requestsq = requestsq.
			Where(
				request.HasGroupWith(
					group.NameEQ(*query.Group),
				),
			)
	}

	requests, err := requestsq.All(ctx)
	if err != nil {
		return nil, err
	}

	reqres := []*RequestResponse{}
	for _, request := range requests {
		reqres = append(reqres, convertEntRequestResponseToModelRequestResponse(request, request.Edges.Tag, request.Edges.Group, request.Edges.Status[0], request.Edges.User))
	}
	return reqres, nil
}

func (repo *EntRepository) CreateRequest(ctx context.Context, amount int, title string, tags []*Tag, targets []*Target, group *Group, userID uuid.UUID) (*RequestDetail, error) {
	tx, txClient, err := setupTx(repo.client, ctx)
	if err != nil {
		return nil, err
	}

	var tagIDs []uuid.UUID
	for _, tag := range tags {
		tagIDs = append(tagIDs, tag.ID)
	}
	created, err := txClient.Request.
		Create().
		SetTitle(title).
		SetAmount(amount).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetUserID(userID).
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}
	user, err := created.QueryUser().Select(user.FieldID).First(ctx)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}
	if group != nil {
		_, err = txClient.Group.
			UpdateOneID(group.ID).
			AddRequest(created).
			Save(ctx)
		if err != nil {
			err = rollBackWithErr(tx, ctx, err)
			return nil, err
		}
	}
	status, err := txClient.RequestStatus.
		Create().
		SetStatus(requeststatus.StatusSubmitted).
		SetCreatedAt(time.Now()).
		SetRequest(created).
		SetUser(user).
		Save(ctx)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}
	modelTargets, err := createRequestTargets(txClient, ctx, created.ID, targets)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("committing transaction: %v", err)
	}
	reqdetail := &RequestDetail{
		ID:        created.ID,
		Status:    convertEntRequestStatusToModelStatus(&status.Status),
		Amount:    created.Amount,
		Title:     created.Title,
		Tags:      tags,
		Targets:   modelTargets,
		Group:     group,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
		CreatedBy: user.ID,
	}
	return reqdetail, nil
}

func (repo *EntRepository) GetRequest(ctx context.Context, requestID uuid.UUID) (*RequestDetail, error) {
	request, err := repo.client.Request.
		Query().
		Where(request.IDEQ(requestID)).
		WithTag().
		WithGroup().
		WithStatus(func(q *ent.RequestStatusQuery) {
			q.Order(ent.Desc(requeststatus.FieldCreatedAt)).Limit(1)
		}).
		WithTarget().
		WithUser().
		Only(ctx)
	if err != nil {
		return nil, err
	}
	var tags []*Tag
	for _, tag := range request.Edges.Tag {
		tags = append(tags, ConvertEntTagToModelTag(tag))
	}
	var targets []*TargetDetail
	for _, target := range request.Edges.Target {
		targets = append(targets, convertEntRequestTargetToModelRequestTarget(target))
	}
	group := ConvertEntGroupToModelGroup(request.Edges.Group)
	reqdetail := &RequestDetail{
		ID:        request.ID,
		Status:    convertEntRequestStatusToModelStatus(&request.Edges.Status[0].Status),
		Amount:    request.Amount,
		Title:     request.Title,
		Tags:      tags,
		Targets:   targets,
		Group:     group,
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
		CreatedBy: request.Edges.User.ID,
	}
	return reqdetail, nil
}

func (repo *EntRepository) UpdateRequest(ctx context.Context, requestID uuid.UUID, amount int, title string, tags []*Tag, targets []*Target, group *Group) (*RequestDetail, error) {
	tx, txClient, err := setupTx(repo.client, ctx)
	if err != nil {
		return nil, err
	}

	var tagIDs []uuid.UUID
	for _, tag := range tags {
		tagIDs = append(tagIDs, tag.ID)
	}
	updated, err := txClient.Request.
		UpdateOneID(requestID).
		SetAmount(amount).
		SetTitle(title).
		ClearTag().
		AddTagIDs(tagIDs...).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}

	if group != nil {
		_, err = txClient.Request.
			UpdateOneID(requestID).
			ClearGroup().
			SetGroupID(group.ID).
			Save(ctx)
	} else {
		_, err = txClient.Request.
			UpdateOneID(requestID).
			ClearGroup().
			Save(ctx)
	}
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}

	status, err := updated.QueryStatus().Select(requeststatus.FieldStatus).First(ctx)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}
	user, err := updated.QueryUser().Select(user.FieldID).First(ctx)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}
	enttags, err := updated.QueryTag().All(ctx)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}
	var modeltags []*Tag
	for _, tag := range enttags {
		modeltags = append(modeltags, ConvertEntTagToModelTag(tag))
	}
	var entgroup *ent.Group
	if group != nil {
		entgroup, err = updated.QueryGroup().Only(ctx)
		if err != nil {
			err = rollBackWithErr(tx, ctx, err)
			return nil, err
		}
	}
	// targets tableに関わることはまとめないとdead lockの原因になる
	err = deleteRequestTargets(txClient, ctx, requestID)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}
	modelTargets, err := createRequestTargets(txClient, ctx, updated.ID, targets)
	if err != nil {
		err = rollBackWithErr(tx, ctx, err)
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("committing transaction: %v", err)
	}
	reqdetail := &RequestDetail{
		ID:        updated.ID,
		Status:    convertEntRequestStatusToModelStatus(&status.Status),
		Amount:    updated.Amount,
		Title:     updated.Title,
		Tags:      modeltags,
		Targets:   modelTargets,
		Group:     ConvertEntGroupToModelGroup(entgroup),
		CreatedAt: updated.CreatedAt,
		UpdatedAt: updated.UpdatedAt,
		CreatedBy: user.ID,
	}
	return reqdetail, nil
}

func convertEntRequestToModelRequest(request *ent.Request) *Request {
	if request == nil {
		return nil
	}
	return &Request{
		ID:        request.ID,
		Amount:    request.Amount,
		CreatedAt: request.CreatedAt,
	}
}

func convertEntRequestResponseToModelRequestResponse(request *ent.Request, tags []*ent.Tag, group *ent.Group, status *ent.RequestStatus, user *ent.User) *RequestResponse {
	if request == nil {
		return nil
	}
	modeltags := []*Tag{}
	for _, tag := range tags {
		modeltags = append(modeltags, ConvertEntTagToModelTag(tag))
	}
	return &RequestResponse{
		ID:        request.ID,
		Status:    convertEntRequestStatusToModelStatus(&status.Status),
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
		CreatedBy: user.ID,
		Amount:    request.Amount,
		Title:     request.Title,
		Tags:      modeltags,
		Group:     ConvertEntGroupToModelGroup(group),
	}
}
