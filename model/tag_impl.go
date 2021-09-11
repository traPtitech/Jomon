package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
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
		modeltags = append(modeltags, convertEntTagToModelTag(tag))
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
	return convertEntTagToModelTag(created), nil
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
	return convertEntTagToModelTag(tag), nil
}

func (repo *EntRepository) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	err := repo.client.Tag.
		DeleteOneID(tagID).
		Exec(ctx)
	return err
}

func convertEntTagToModelTag(tag *ent.Tag) *Tag {
	return &Tag{
		ID:          tag.ID,
		Name:        tag.Name,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
		DeletedAt:   tag.DeletedAt,
	}
}
