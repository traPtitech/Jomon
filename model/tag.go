package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/tag"
)

type TagRepository interface {
	GetTags(ctx context.Context, tagID uuid.UUID) (*ent.Tag, error)
	UpdateTag(ctx context.Context, tagID uuid.UUID, name string, description string) (*ent.Tag, error)
	GetTagTransactions(ctx context.Context, tag *ent.Tag) ([]*ent.Transaction, error)
	GetTagRequests(ctx context.Context, tag *ent.Tag) ([]*ent.Request, error)
}

func (repo *EntRepository) GetTags(ctx context.Context, tagID uuid.UUID) (*ent.Tag, error) {
	tag, err := repo.client.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (repo *EntRepository) UpdateTag(ctx context.Context, tagID uuid.UUID, name string, description string) (*ent.Tag, error) {
	tag, err := repo.client.Tag.
		UpdateOneID(tagID).
		SetName(name).
		SetDescription(description).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func (repo *EntRepository) GetTagTransactions(ctx context.Context, tag *ent.Tag) ([]*ent.Transaction, error) {
	transactions, err := tag.
		QueryTransaction().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (repo *EntRepository) GetTagRequests(ctx context.Context, tag *ent.Tag) ([]*ent.Request, error) {
	requests, err := tag.
		QueryRequest().
		All(ctx)
	if err != nil {
		return nil, err
	}
	return requests, nil
}
