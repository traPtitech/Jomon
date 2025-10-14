package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/ent/tag"
	"github.com/traPtitech/Jomon/internal/service"
)

func (repo *EntRepository) GetTags(ctx context.Context) ([]*Tag, error) {
	tags, err := repo.client.Tag.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modeltags := lo.Map(tags, func(t *ent.Tag, _ int) *Tag {
		return ConvertEntTagToModelTag(t)
	})

	return modeltags, nil
}

func (repo *EntRepository) GetTag(ctx context.Context, tagID uuid.UUID) (*Tag, error) {
	t, err := repo.client.Tag.
		Query().
		Where(tag.IDEQ(tagID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntTagToModelTag(t), nil
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

func (repo *EntRepository) UpdateTag(
	ctx context.Context, tagID uuid.UUID, name string,
) (*Tag, error) {
	t, err := repo.client.Tag.
		UpdateOneID(tagID).
		SetName(name).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntTagToModelTag(t), nil
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
		DeletedAt: service.TimeToNullTime(enttag.DeletedAt).Time,
	}
}
