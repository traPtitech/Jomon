package model

import (
	"context"
	// "errors"

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
	members, err := repo.client.Group.
		Query().
		Where(group.IDEQ(groupID)).
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
		Update().
		Where(group.IDEQ(groupID)).
		AddUserIDs(userID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// var groups []*Group
	// groups, err = repo.GetGroups(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// groupIDExist := false
	// for _, v := range groups {
	// 	if v.ID == groupID {
	// 		groupIDExist = true
	// 	}
	// }
	// if groupIDExist == false {
	// 	err = errors.New("存在しないgroupIDです。")
	// 	return nil, err
	// }

	created := &Member{userID}
	return created, nil
}

func (repo *EntRepository) DeleteMember(ctx context.Context, groupID uuid.UUID, userID uuid.UUID) error {
	_, err := repo.client.Group.
		Update().
		Where(group.IDEQ(groupID)).
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
