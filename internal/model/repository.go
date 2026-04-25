package model

import (
	"github.com/traPtitech/Jomon/internal/ent"
	"github.com/traPtitech/Jomon/internal/service"
)

type EntRepository struct {
	client *ent.Client
}

var _ service.Repository = (*EntRepository)(nil)

func NewEntRepository(client *ent.Client) *EntRepository {
	repo := &EntRepository{
		client: client,
	}
	return repo
}
