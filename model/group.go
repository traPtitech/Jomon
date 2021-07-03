//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE -package=mock_$GOPACKAGE
package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/traPtitech/Jomon/ent/group"
)

type GroupRepository interface {
	GetGroups(ctx context.Context) ([]*Group, error)
	GetGroup(ctx context.Context, groupID uuid.UUID) (*Group, error)
	CreateGroup(ctx context.Context, name string, description string, budget *int, owners *[]User) (*Group, error)
	GetMembers(ctx context.Context, groupID string) ([]*Group, error)
}

type Group struct {
	ID          uuid.UUID
	Name        string
	Description string
	Budget      *int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (repo *EntRepository) GetMembers(ctx context.Context, groupID string) ([]*Group, error) {
	members, err := repo.client.Group.Where(group.IDEQ(groupID))
	if err != nil {
		return nil, err
	}
	modelmembers := []*Group{}
	for _, member := range members {
		modelmembers = append(modelmembers, ConvertEntGroupToModelGroup(member))
	}
	return modelmembers, nil
}
