package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/group"
)

func (repo *EntRepository) GetGroups(ctx context.Context) ([]*Group, error) {
	groups, err := repo.client.Group.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modelgroups := []*Group{}
	for _, group := range groups {
		modelgroups = append(modelgroups, ConvertEntGroupToModelGroup(group))
	}
	return modelgroups, nil
}

func (repo *EntRepository) GetGroup(ctx context.Context, groupID uuid.UUID) (*Group, error) {
	group, err := repo.client.Group.
		Query().
		Where(group.IDEQ(groupID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntGroupToModelGroup(group), nil
}

func (repo *EntRepository) CreateGroup(ctx context.Context, name string, description string, budget *int, owners *[]User) (*Group, error) {
	created, err := repo.client.Group.
		Create().
		SetName(name).
		SetDescription(description).
		SetBudget(*budget).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntGroupToModelGroup(created), nil
}

func ConvertEntGroupToModelGroup(entgroup *ent.Group) *Group {
	if entgroup == nil {
		return nil
	}
	return &Group{
		ID:          entgroup.ID,
		Name:        entgroup.Name,
		Description: entgroup.Description,
		Budget:      entgroup.Budget,
		CreatedAt:   entgroup.CreatedAt,
		UpdatedAt:   entgroup.UpdatedAt,
		DeletedAt:   entgroup.DeletedAt,
	}
}
