package model

import (
	"context"

	"github.com/traPtitech/Jomon/ent"
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
		modelgroups = append(modelgroups, convertEntGroupToModelGroup(group))
	}
	return modelgroups, nil
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
	return convertEntGroupToModelGroup(created), nil
}

func convertEntGroupToModelGroup(group *ent.Group) *Group {
	return &Group{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Budget:      group.Budget,
		CreatedAt:   group.CreatedAt,
		UpdatedAt:   group.UpdatedAt,
		DeletedAt:   group.DeletedAt,
	}
}
