package data

import (
	"context"

	"transaction/internal/biz"
	"transaction/internal/model"

	"gorm.io/gorm"
)

type TransactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepo(data *Data) biz.TransactionRepo {
	return &TransactionRepo{
		db: data.DB,
	}
}


func (r *TransactionRepo) Create(ctx context.Context,transaction *model.Transaction,) error {

	if transaction.RelatedTransactionID == 0 {
		return r.db.WithContext(ctx).
			Omit("related_transaction_id").
			Create(transaction).Error
	}

	return r.db.WithContext(ctx).
		Create(transaction).Error
}


func (r *TransactionRepo) GetByID(ctx context.Context,id int64,) (*model.Transaction, error) {

	var transaction model.Transaction

	err := r.db.WithContext(ctx).
		Where("transaction_id = ?", id).
		First(&transaction).Error

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}


func (r *TransactionRepo) GetAll(ctx context.Context,) ([]*model.Transaction, error) {

	var transactions []*model.Transaction

	err := r.db.WithContext(ctx).
		Find(&transactions).Error

	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepo) Update(ctx context.Context,transaction *model.Transaction,) error {

	return r.db.WithContext(ctx).
		Model(&model.Transaction{}).
		Where(
			"transaction_id = ?",
			transaction.TransactionID,
		).
		Updates(map[string]interface{}{
			"transaction_type": transaction.TransactionType,
			"amount":           transaction.Amount,
			"currency":         transaction.Currency,
			"status":           transaction.Status,
			"description":      transaction.Description,
			"merchant_info":    transaction.MerchantInfo,
			"posting_date":     transaction.PostingDate,
			"posted_at":        transaction.PostedAt,
			"archived_at":      transaction.ArchivedAt,
		}).Error
}

func (r *TransactionRepo) Delete(ctx context.Context,id int64,) error {
	return r.db.WithContext(ctx).
		Delete(&model.Transaction{}, id).Error
}