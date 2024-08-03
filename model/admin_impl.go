package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/ent"
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

	admins := lo.Map(users, func(u *ent.User, insex int) *Admin {
		return &Admin{
			ID: u.ID,
		}
	})

	return admins, nil
}

func (repo *EntRepository) AddAdmins(ctx context.Context, userIDs []uuid.UUID) error {
	_, err := repo.client.User.
		Update().
		Where(user.IDIn(userIDs...)).
		SetAdmin(true).
		Save(ctx)

	return err
}

func (repo *EntRepository) DeleteAdmins(ctx context.Context, userIDs []uuid.UUID) error {
	_, err := repo.client.User.
		Update().
		Where(user.IDIn(userIDs...)).
		SetAdmin(false).
		Save(ctx)

	return err
}
