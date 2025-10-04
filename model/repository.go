package model

import (
	"github.com/traPtitech/Jomon/ent"
)

type Repository interface {
	AccountManagerRepository
	CommentRepository
	FileRepository
	ApplicationFileRepository
	ApplicationStatusRepository
	ApplicationTagRepository
	ApplicationTargetRepository
	ApplicationRepository
	TagRepository
	UserRepository
}

type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client) *EntRepository {
	repo := &EntRepository{
		client: client,
	}
	return repo
}
