package model

import (
	"context"

	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/user"
)

func (repo *EntRepository) GetMe(ctx context.Context, name string) (*User, error) {
	user, err := repo.client.User.
		Query().
		Where(user.NameEQ(name)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntUserToModelUser(user), nil
}

func ConvertEntUserToModelUser(user *ent.User) *User {
	return &User{
		ID:          user.ID,
		Name:        user.Name,
		DisplayName: user.DisplayName,
		Admin:       user.Admin,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		DeletedAt:   user.DeletedAt,
	}
}
