package model

import (
	"context"
	"errors"

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
	gotGroup, err := repo.client.Group.
		Query().
		Where(group.IDEQ(groupID)).
		First(ctx)
	if err != nil {
		return nil, err
	}
	if gotGroup == nil {
		return nil, errors.New("unknown group id")
	}

	members, err := gotGroup.
		QueryUser().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modelmembers := []*User{}
	for _, member := range members {
		modelmembers = append(modelmembers, ConvertEntUserToModelUser(member))
	}
	return modelmembers, nil
}

func (repo *EntRepository) CreateMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) (*Member, error) {
	_, err := repo.client.Group.
		UpdateOneID(groupID).
		AddUserIDs(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}
<<<<<<< HEAD
=======

>>>>>>> 894e3e324c7bd810f1f210fbf313cb97e8038569
	created := &Member{userID}
	return created, nil
}

func (repo *EntRepository) DeleteMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) error {
	gotUser, err := repo.client.User.
		Query().
		Where(user.IDEQ(userID)).
		First(ctx)
	if err != nil {
		return err
	}
	if gotUser == nil {
		return errors.New("unknown user id")
	}
	
	_, err = repo.client.Group.
		UpdateOneID(groupID).
		RemoveUserIDs(userID).
		Save(ctx)
	if err != nil {
		return err
	}
	return nil
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
