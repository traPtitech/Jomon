package model

import (
	"context"

	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/group"

	"github.com/google/uuid"
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

func (repo *EntRepository) GetOwners(ctx context.Context, GroupID uuid.UUID) ([]*Owner, error) {
	groupowners, err := repo.client.Group.
		Query().
		Where(group.IDEQ(GroupID)).
		QueryOwner().
		All(ctx)
	if err != nil {
		return nil, err
	}
	owners := []*Owner{}
	for _, groupowner := range groupowners {
		owners = append(owners, &Owner{ID: groupowner.ID})
	}

	return owners, nil
}

func (repo *EntRepository) CreateOwner(ctx context.Context, GroupID uuid.UUID, OwnerID uuid.UUID) (*Owner, error) {
	_, err := repo.client.Group.
		Update().
		Where(group.IDEQ(GroupID)).
		AddOwnerIDs(OwnerID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	resowner := &Owner{
		ID: OwnerID,
	}
	return resowner, nil

}

func (repo *EntRepository) DeleteOwner(ctx context.Context, GroupID uuid.UUID, OwnerID uuid.UUID) error {
	_, err := repo.client.Group.
		Update().
		Where(group.IDEQ(GroupID)).
		RemoveOwnerIDs(OwnerID).
		Save(ctx)

	return err
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
