package model

import (
	"github.com/traPtitech/Jomon/ent"
	storagePkg "github.com/traPtitech/Jomon/storage"
)

type Repository interface {
	AdminRepository
	CommentRepository
	FileRepository
	GroupBudgetRepository
	GroupRepository
	RequestFileRepository
	RequestStatusRepository
	RequestTagRepository
	RequestTargetRepository
	RequestRepository
	TagRepository
	TransactionDetailRepository
	TransactionTagRepository
	TransactionRepository
	UserRepository
}

type EntRepository struct {
	client *ent.Client
}

func NewEntRepository(client *ent.Client, storage storagePkg.Storage) *EntRepository {
	repo := &EntRepository{
		client: client,
	}
	return repo
}
