package model

import (
	"github.com/traPtitech/Jomon/ent"
)

type Repository interface {
	CommentRepository
	FileRepository
	GroupBudgetRepository
	GroupOwnerRepository
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

func NewEntRepository(client *ent.Client) Repository {
	repo := &EntRepository{
		client: client,
	}
	return repo
}
