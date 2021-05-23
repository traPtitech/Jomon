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
		modelgroups = append(modelgroups, ConvertEntGroupToModelGroup(group))
	}
	return modelgroups, nil
}

func ConvertEntGroupToModelGroup(entgroup *ent.Group) *Group {
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
