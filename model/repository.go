package model

import (
	"github.com/traPtitech/Jomon/ent"
	storagePkg "github.com/traPtitech/Jomon/storage"
)

type Repository interface {
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
	client  *ent.Client
	storage storagePkg.Storage
}

func NewEntRepository(client *ent.Client, storage storagePkg.Storage) Repository {
	repo := &EntRepository{
		client:  client,
		storage: storage,
	}
	return repo
}
