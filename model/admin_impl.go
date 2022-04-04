package model

import (
	"context"

	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/user"
	"google.golang.org/appengine/user"
)

func (repo *EntRepository) GetAdmins(ctx context.Context) ([]*Admin, error) {
	admins, err := repo.client.User.
		Query().
		Where(user.Admin(true)).
		All(ctx)
	if err != nil {
		return nil, err
	}

	return ConvertEntUserToModelAdmin(admins), nil
}

func ConvertEntUserToModelAdmin(entUser *ent.User) *Admin {
	return &Admin{
		ID: entUser.ID,
	}
}
