package model

import (
	"time"

	"github.com/jinzhu/gorm"

	"github.com/gofrs/uuid"
)

// Transaction struct of Transaction
type Transaction struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	Amount    int       `gorm:"type:int(11);not null" json:"amount"`
	Target    string    `gorm:"type:varchar(64);not null" json:"target"`
	RequestID uuid.UUID `gorm:"char(36);index" json:"request_id`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP;index" json:"created_at"`
}

// TransactionRepository トランザクションのリポジトリ
type TransactionRepository interface {
	GetTransaction(id uuid.UUID) (Transaction, error)
	GetTransactionList(sort string,
		currentStatus *Status,
		financialYear *int,
		applicant string,
		submittedSince *time.Time,
		submittedUntil *time.Time,
	) ([]Request, error)
	CreateTransaction(
		amount int,
		target string,
		requestID *uuid.UUID,
	) (uuid.UUID, error)
	PatchTransaction(
		trnsID uuid.UUID,
		amount *int,
		target string,
		requestID *uuid.UUID,
	) error
}

type transactionRepository struct{}

// NewTransactionRepository Make TransactionRepository
func NewTransactionRepository() RequestRepository {
	return &requestRepository{}
}

func (*transactionRepository) GetTransaction(id uuid.UUID) (Transaction, error) {
	var trns Transaction

	err := db.Order("created_at").Find(&trns, Request{ID: id}).Error
	if err != nil {
		return Transaction{}, err
	}

	return trns, nil
}

func (*transactionRepository) GetTransactionList(sort string, financialYear *int, applicant string, submittedSince *time.Time, submittedUntil *time.Time) ([]Transaction, error) {
	query := db

	if financialYear != nil {
		financialYear := time.Date(*financialYear, 4, 1, 0, 0, 0, 0, time.Local)
		financialYearEnd := financialYear.AddDate(1, 0, 0)
		query = query.Where("created_at >= ?", financialYear).Where("created_at < ?", financialYearEnd)
	}

	if applicant != "" {
		query = query.Where("created_by = ?", applicant)
	}

	if submittedSince != nil {
		query = query.Where("created_at >= ?", *submittedSince)
	}

	if submittedUntil != nil {
		query = query.Where("created_at < ?", *submittedUntil)
	}

	switch sort {
	case "", "created_at":
		query = query.Order("created_at desc")
	case "-created_at":
		query = query.Order("created_at")
	case "title":
		query = query.Order("title")
	case "-title":
		query = query.Order("title desc")
	}

	//noinspection GoPreferNilSlice
	apps := []Transaction{}
	err := query.Find(&apps).Error
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (repo *transactionRepository) CreateTransaction(amount int, targets []string, requestID *uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID
	id, err := uuid.NewV4()
	if err != nil {
		return uuid.Nil, err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		for _, target := range targets {
			if err = repo.createTransaction(tx, id, amount, target, *requestID); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (repo *transactionRepository) PatchTransaction(trnsID uuid.UUID, amount *int, transactionTargets []string, requestID *uuid.UUID) error {
	return db.Transaction(func(tx *gorm.DB) error {
		trns, err := repo.GetTransaction(trnsID)
		if err != nil {
			return err
		}

		trns.ID, err = uuid.NewV4() // zero value of int is 0
		if err != nil {
			return err
		}

		if amount != nil {
			trns.Amount = *amount
		}

		if len(transactionTargets) > 0 {
			for _, userID := range transactionTargets {
				if err = repo.createTransactionTarget(tx, trnsID, userID); err != nil {
					return err
				}
			}
		}

		return tx.Model(&Transaction{ID: trnsID}).Updates(Transaction{}).Error
	})
}

func (*transactionRepository) createTransaction(db *gorm.DB, id uuid.UUID, amount int, target string, requestID uuid.UUID) error {
	trns := Transaction{
		ID:        id,
		Amount:    amount,
		Target:    target,
		RequestID: requestID,
	}

	err := db.Create(&trns).Error
	if err != nil {
		return err
	}

	return nil
}
