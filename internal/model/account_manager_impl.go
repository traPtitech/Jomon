package model

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/ent/user"
)

func (repo *EntRepository) GetAccountManagers(ctx context.Context) ([]*AccountManager, error) {
	users, err := repo.client.User.
		Query().
		Where(user.AccountManager(true)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	accountManagers := lo.Map(users, func(u *ent.User, _ int) *AccountManager {
		return &AccountManager{
			ID: u.ID,
		}
	})

	return accountManagers, nil
}

func (repo *EntRepository) AddAccountManagers(ctx context.Context, userIDs []uuid.UUID) error {
	_, err := repo.client.User.
		Update().
		Where(user.IDIn(userIDs...)).
		SetAccountManager(true).
		Save(ctx)

	return err
}

func (repo *EntRepository) DeleteAccountManagers(ctx context.Context, userIDs []uuid.UUID) error {
	_, err := repo.client.User.
		Update().
		Where(user.IDIn(userIDs...)).
		SetAccountManager(false).
		Save(ctx)

	return err
}
