package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/group"
	"github.com/traPtitech/Jomon/ent/user"
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

func (repo *EntRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]*User, error) {
	members, err := repo.client.Group.Query().Where(group.IDEQ(groupID)).QueryUser().All(ctx)
	if err != nil {
		return nil, err
	}
	modelmembers := []*User{}
	for _, member := range members {
		modelmembers = append(modelmembers, ConvertEntUserToModelUser(member))
	}
	return modelmembers, nil
}

func (repo *EntRepository) PostMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (*User, error) {
	group, err := repo.client.Group.Query().Where(group.IDEQ(groupID)).First(ctx)
	if err != nil {
		return nil, err
	}
	_, err = group.Update().AddUserIDs(userID).Save(ctx)
	if err != nil {
		return nil, err
	}
	created, err := repo.client.User.Query().Where(user.IDEQ(userID)).First(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntUserToModelUser(created), nil
}

func (repo *EntRepository) DeleteMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) error {
	group, err := repo.client.Group.Query().Where(group.IDEQ(groupID)).First(ctx)
	if err != nil {
		return err
	}

	_, err = group.Update().RemoveUserIDs(userID).Save(ctx)
	if err != nil {
		return err
	}

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
