package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/request"
	"github.com/traPtitech/Jomon/ent/requeststatus"
	"github.com/traPtitech/Jomon/ent/user"
)

func (repo *EntRepository) GetRequests(ctx context.Context, query RequestQuery) ([]*RequestResponse, error) {
	// Querying
	var requests []*ent.Request
	var err error
	if query.Sort == nil || *query.Sort == "created_at" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt)).Limit(1)
			}).
			WithUser().
			Order(ent.Desc(request.FieldCreatedAt)).
			All(ctx)
	} else if *query.Sort == "-created_at" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt)).Limit(1)
			}).
			WithUser().
			Order(ent.Asc(request.FieldCreatedAt)).
			All(ctx)
	} else if *query.Sort == "title" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt)).Limit(1)
			}).
			WithUser().
			Order(ent.Desc(request.FieldTitle)).
			All(ctx)
	} else if *query.Sort == "-title" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus(func(q *ent.RequestStatusQuery) {
				q.Order(ent.Desc(requeststatus.FieldCreatedAt)).Limit(1)
			}).
			WithUser().
			Order(ent.Asc(request.FieldTitle)).
			All(ctx)
	}
	if err != nil {
		return nil, err
	}

	var reqres []*RequestResponse
	for _, request := range requests {
		reqres = append(reqres, ConvertEntRequestResponseToModelRequestResponse(request, request.Edges.Tag, request.Edges.Group, request.Edges.Status[0], request.Edges.User))
	}
	return reqres, nil
}

func (repo *EntRepository) CreateRequest(ctx context.Context, amount int, title string, content string, tags []*Tag, group *Group, userID uuid.UUID) (*RequestDetail, error) {
	// TODO: WIP
	var tagIDs []uuid.UUID
	for _, tag := range tags {
		tagIDs = append(tagIDs, tag.ID)
	}
	created, err := repo.client.Request.
		Create().
		SetTitle(title).
		SetAmount(amount).
		SetContent(content).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		SetUserID(userID).
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	user, err := created.QueryUser().Select(user.FieldID).First(ctx)
	if err != nil {
		return nil, err
	}
	if group != nil {
		_, err = repo.client.Group.
			UpdateOneID(group.ID).
			AddRequest(created).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}
	status, err := repo.client.RequestStatus.
		Create().
		SetStatus(requeststatus.StatusSubmitted).
		SetReason("").
		SetCreatedAt(time.Now()).
		SetRequest(created).
		SetUser(user).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	reqdetail := &RequestDetail{
		ID:        created.ID,
		Status:    string(status.Status),
		Amount:    created.Amount,
		Title:     created.Title,
		Content:   created.Content,
		Tags:      tags,
		Group:     group,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
		CreatedBy: user.ID,
	}
	return reqdetail, nil
}

func ConvertEntRequestToModelRequest(request *ent.Request) *Request {
	if request == nil {
		return nil
	}
	return &Request{
		ID:        request.ID,
		Amount:    request.Amount,
		CreatedAt: request.CreatedAt,
	}
}

func ConvertEntRequestResponseToModelRequestResponse(request *ent.Request, tags []*ent.Tag, group *ent.Group, status *ent.RequestStatus, user *ent.User) *RequestResponse {
	if request == nil {
		return nil
	}
	var modeltags []*Tag
	for _, tag := range tags {
		modeltags = append(modeltags, ConvertEntTagToModelTag(tag))
	}
	return &RequestResponse{
		ID:        request.ID,
		Status:    string(status.Status),
		CreatedAt: request.CreatedAt,
		UpdatedAt: request.UpdatedAt,
		CreatedBy: user.ID,
		Amount:    request.Amount,
		Title:     request.Title,
		Content:   request.Content,
		Tags:      modeltags,
		Group:     ConvertEntGroupToModelGroup(group),
	}
}
