package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/group"
	"github.com/traPtitech/Jomon/ent/user"
	"github.com/traPtitech/Jomon/service"
)

func (repo *EntRepository) GetGroups(ctx context.Context) ([]*Group, error) {
	groups, err := repo.client.Group.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	modelgroups := lo.Map(groups, func(g *ent.Group, _ int) *Group {
		return ConvertEntGroupToModelGroup(g)
	})
	return modelgroups, nil
}

func (repo *EntRepository) GetGroup(ctx context.Context, groupID uuid.UUID) (*Group, error) {
	g, err := repo.client.Group.
		Query().
		Where(group.IDEQ(groupID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntGroupToModelGroup(g), nil
}

func (repo *EntRepository) CreateGroup(
	ctx context.Context, name string, description string, budget *int,
) (*Group, error) {
	created, err := repo.client.Group.
		Create().
		SetName(name).
		SetDescription(description).
		SetNillableBudget(budget).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return ConvertEntGroupToModelGroup(created), nil
}

func (repo *EntRepository) UpdateGroup(
	ctx context.Context, groupID uuid.UUID, name string, description string, budget *int,
) (*Group, error) {
	updated, err := repo.client.Group.
		UpdateOneID(groupID).
		SetName(name).
		SetDescription(description).
		ClearBudget().
		SetNillableBudget(budget).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return ConvertEntGroupToModelGroup(updated), nil
}

func (repo *EntRepository) DeleteGroup(ctx context.Context, groupID uuid.UUID) error {
	err := repo.client.Group.
		DeleteOneID(groupID).
		Exec(ctx)

	return err
}

func (repo *EntRepository) GetOwners(ctx context.Context, groupID uuid.UUID) ([]*Owner, error) {
	groupowners, err := repo.client.Group.
		Query().
		Where(group.IDEQ(groupID)).
		QueryOwner().
		Select(user.FieldID).
		All(ctx)
	if err != nil {
		return nil, err
	}
	owners := lo.Map(groupowners, func(groupowner *ent.User, _ int) *Owner {
		return &Owner{ID: groupowner.ID}
	})

	return owners, nil
}

func (repo *EntRepository) AddOwners(
	ctx context.Context, groupID uuid.UUID, ownerIDs []uuid.UUID,
) ([]*Owner, error) {
	_, err := repo.client.Group.
		Update().
		Where(group.IDEQ(groupID)).
		AddOwnerIDs(ownerIDs...).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	resowners := lo.Map(ownerIDs, func(owner uuid.UUID, _ int) *Owner {
		return &Owner{ID: owner}
	})

	return resowners, nil
}

func (repo *EntRepository) DeleteOwners(
	ctx context.Context, groupID uuid.UUID, ownerIDs []uuid.UUID,
) error {
	_, err := repo.client.Group.
		Update().
		Where(group.IDEQ(groupID)).
		RemoveOwnerIDs(ownerIDs...).
		Save(ctx)

	return err
}

func (repo *EntRepository) GetMembers(ctx context.Context, groupID uuid.UUID) ([]*Member, error) {
	members, err := repo.client.Group.
		Query().
		Where(group.IDEQ(groupID)).
		QueryUser().
		Select(user.FieldID).
		All(ctx)

	if err != nil {
		return nil, err
	}
	modelmembers := lo.Map(members, func(member *ent.User, _ int) *Member {
		return &Member{member.ID}
	})

	return modelmembers, nil
}

func (repo *EntRepository) AddMembers(
	ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID,
) ([]*Member, error) {
	_, err := repo.client.Group.
		Update().
		Where(group.IDEQ(groupID)).
		AddUserIDs(userIDs...).
		Save(ctx)

	if err != nil {
		return nil, err
	}
	resMembers := lo.Map(userIDs, func(member uuid.UUID, _ int) *Member {
		return &Member{member}
	})

	return resMembers, nil
}

func (repo *EntRepository) DeleteMembers(
	ctx context.Context, groupID uuid.UUID, userIDs []uuid.UUID,
) error {
	_, err := repo.client.Group.
		Update().
		Where(group.IDEQ(groupID)).
		RemoveUserIDs(userIDs...).
		Save(ctx)

	if err != nil {
		return err
	}
	return nil
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
		DeletedAt:   service.TimeToNullTime(entgroup.DeletedAt),
	}
}
