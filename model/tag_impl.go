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

func (repo *EntRepository) CreateTag(ctx context.Context, name string) (*Tag, error) {
	created, err := repo.client.Tag.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntTagToModelTag(created), nil
}

func (repo *EntRepository) UpdateTag(ctx context.Context, tagID uuid.UUID, name string) (*Tag, error) {
	tag, err := repo.client.Tag.
		UpdateOneID(tagID).
		SetName(name).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntTagToModelTag(tag), nil
}

func (repo *EntRepository) DeleteTag(ctx context.Context, tagID uuid.UUID) error {
	err := repo.client.Tag.
		DeleteOneID(tagID).
		Exec(ctx)
	return err
}

func ConvertEntTagToModelTag(enttag *ent.Tag) *Tag {
	return &Tag{
		ID:        enttag.ID,
		Name:      enttag.Name,
		CreatedAt: enttag.CreatedAt,
		UpdatedAt: enttag.UpdatedAt,
		DeletedAt: enttag.DeletedAt,
	}
}
