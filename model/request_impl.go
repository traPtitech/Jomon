package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/request"
)

func (repo *EntRepository) GetRequests(ctx context.Context, query RequestQuery) ([]*RequestResponse, error) {
	// Querying
	var requests []*ent.Request
	var err error
	if *query.Sort == "created_at" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus().
			WithUser().
			Order(ent.Desc(request.FieldCreatedAt)).
			All(ctx)
	} else if *query.Sort == "-created_at" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus().
			WithUser().
			Order(ent.Asc(request.FieldCreatedAt)).
			All(ctx)
	} else if *query.Sort == "title" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus().
			WithUser().
			Order(ent.Desc(request.FieldTitle)).
			All(ctx)
	} else if *query.Sort == "-title" {
		requests, err = repo.client.Request.
			Query().
			WithTag().
			WithGroup().
			WithStatus().
			WithUser().
			Order(ent.Asc(request.FieldTitle)).
			All(ctx)
	}
	if err != nil {
		return nil, err
	}

	var reqres []*RequestResponse
	for _, request := range requests {
		reqres = append(reqres, ConvertEntRequestResponseToModelRequestResponse(request, request.Edges.Tag, request.Edges.Group, request.Edges.User.Edges.RequestStatus, request.Edges.User))
	}
	return reqres, nil
}

func (repo *EntRepository) CreateRequest(ctx context.Context, amount int, title string, content string, tags []*Tag, group Group, files []*File) (*RequestDetail, error) {
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
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	_, err = repo.client.Group.
		UpdateOneID(group.ID).
		AddRequest(created).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	reqdetail := &RequestDetail{
		ID:      created.ID,
		Amount:  created.Amount,
		Title:   created.Title,
		Content: created.Content,
		Tags:    tags,
		Group:   group,
	}
	return reqdetail, nil
}

func ConvertEntRequestToModelRequest(request *ent.Request) *Request {
	return &Request{
		ID:        request.ID,
		Amount:    request.Amount,
		CreatedAt: request.CreatedAt,
	}
}

func ConvertEntRequestResponseToModelRequestResponse(request *ent.Request, tags []*ent.Tag, group *ent.Group, status *ent.RequestStatus, user *ent.User) *RequestResponse {
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
		Group:     *ConvertEntGroupToModelGroup(group),
	}
}
