package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent"
	"github.com/traPtitech/Jomon/ent/user"
)

func (repo *EntRepository) GetUsers(ctx context.Context) ([]*User, error) {
	users, err := repo.client.User.
		Query().
		All(ctx)
	if err != nil {
		return nil, err
	}
	var modelusers []*User
	for _, user := range users {
		modelusers = append(modelusers, ConvertEntUserToModelUser(user))
	}
	return modelusers, nil
}

func (repo *EntRepository) CreateUser(ctx context.Context, name string, dn string, admin bool) (*User, error) {
	user, err := repo.client.User.
		Create().
		SetName(name).
		SetDisplayName(dn).
		SetAdmin(admin).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntUserToModelUser(user), nil
}

func (repo *EntRepository) GetUserByName(ctx context.Context, name string) (*User, error) {
	user, err := repo.client.User.
		Query().
		Where(user.NameEQ(name)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntUserToModelUser(user), nil
}

func (repo *EntRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	user, err := repo.client.User.
		Query().
		Where(user.IDEQ(userID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntUserToModelUser(user), nil
}

func (repo *EntRepository) UpdateUser(ctx context.Context, userID uuid.UUID, name string, dn string, admin bool) (*User, error) {
	user, err := repo.client.User.
		UpdateOneID(userID).
		SetName(name).
		SetDisplayName(dn).
		SetAdmin(admin).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return ConvertEntUserToModelUser(user), nil
}

func ConvertEntUserToModelUser(user *ent.User) *User {
	if user == nil {
		return nil
	}
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
