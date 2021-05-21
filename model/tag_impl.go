package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/tag"
)

func (repo *EntRepository) GetTags(ctx context.Context) ([]*Tag, error) {
	tags, err := repo.client.Tag.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modeltags := []*Tag{}
	for _, tag := range tags {
		modeltags = append(modeltags, ConvertEntTagToModelTag(tag))
	}
	return modeltags, nil
}

func (repo *EntRepository) CreateTag(ctx context.Context, name string, description string) (*Tag, error) {
	created, err := repo.client.Tag.
		Create().
		SetName(name).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntTagToModelTag(created), nil
}

func (repo *EntRepository) GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error) {
	tag, err := repo.client.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntTagToModelTag(tag), nil
}

func (repo *EntRepository) UpdateTag(ctx context.Context, tagID uuid.UUID, name string, description string) (*Tag, error) {
	tag, err := repo.client.Tag.
		UpdateOneID(tagID).
		SetName(name).
		SetDescription(description).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntTagToModelTag(tag), nil
}

func (repo *EntRepository) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	tag, err := repo.client.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return err
	}
	err = repo.client.Tag.
		DeleteOne(tag).
		Exec(ctx)
	return err
}

func (repo *EntRepository) GetTagTransactions(ctx context.Context, tagID uuid.UUID) ([]*Transaction, error) {
	tag, err := repo.client.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	transactions, err := tag.
		QueryTransaction().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modeltransactions := []*Transaction{}
	for _, transaction := range transactions {
		modeltransactions = append(modeltransactions, ConvertEntTransactionToModelTransaction(transaction))
	}
	return modeltransactions, nil
}

func (repo *EntRepository) GetTagRequests(ctx context.Context, tagID uuid.UUID) ([]*Request, error) {
	tag, err := repo.client.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	requests, err := tag.
		QueryRequest().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modelrequests := []*Request{}
	for _, request := range requests {
		modelrequests = append(modelrequests, ConvertEntRequestToModelRequest(request))
	}
	return modelrequests, nil
}

func ConvertEntTagToModelTag(enttag *ent.Tag) *Tag {
	return &Tag{
		ID:          enttag.ID,
		Name:        enttag.Name,
		Description: enttag.Description,
		CreatedAt:   enttag.CreatedAt,
		UpdatedAt:   enttag.UpdatedAt,
		DeletedAt:   enttag.DeletedAt,
	}
}
