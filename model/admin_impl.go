package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/user"
)

func (repo *EntRepository) GetAdmins(ctx context.Context) ([]*Admin, error) {
	users, err := repo.client.User.
		Query().
		Where(user.Admin(true)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	admins := []*Admin{}
	for _, user := range users {
		admins = append(admins, &Admin{
			ID: user.ID,
		})
	}

	return admins, nil
}

func (repo *EntRepository) CreateAdmin(ctx context.Context, userID uuid.UUID) (*Admin, error) {
	user, err := repo.client.User.
		Query().
		Where(user.ID(userID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	if user.Admin {
		return nil, nil
	}

	_, err = repo.client.User.
		UpdateOne(user).
		SetAdmin(true).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &Admin{
		ID: user.ID,
	}, nil
}

func (repo *EntRepository) DeleteAdmin(ctx context.Context, userID uuid.UUID) error {
	user, err := repo.client.User.
		Query().
		Where(user.ID(userID)).
		Only(ctx)
	if err != nil {
		return err
	}

	if !user.Admin {
		return nil
	}

	_, err = repo.client.User.
		UpdateOne(user).
		SetAdmin(false).
		Save(ctx)
	if err != nil {
		return err
	}

	return nil
}
