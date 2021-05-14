package model

import (
	"gorm.io/gorm"
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
}

type GormRepository struct {
	db *gorm.DB
}
