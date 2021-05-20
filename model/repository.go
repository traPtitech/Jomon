package model

import (
	"github.com/traPtitech/Jomon/ent"
)

type Repository interface {
	AdministratorRepository
	CommentRepository
	FileRepository
	GroupBudgetRepository
	GroupOwnerRepository
	GroupUserRepository
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
