package data

import (
	"context"

	"payments/internal/biz"
	"payments/internal/model"

	"gorm.io/gorm"
)

type PaymentRepo struct {
	db *gorm.DB
}

func NewPaymentRepo(data *Data) biz.PaymentRepo {
	return &PaymentRepo{
		db: data.DB,
	}
}

func (r *PaymentRepo) Create(ctx context.Context,payment *model.Payment,
) error {

	return r.db.WithContext(ctx).
		Create(payment).
		Error
}

func (r *PaymentRepo) GetByID(ctx context.Context,id int64,) (*model.Payment, error) {

	var payment model.Payment

	err := r.db.WithContext(ctx).
		Where("payment_id = ?", id).
		First(&payment).
		Error

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentRepo) GetAll(ctx context.Context,) ([]*model.Payment, error) {

	var payments []*model.Payment

	err := r.db.WithContext(ctx).
		Find(&payments).
		Error

	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepo) Update(ctx context.Context,payment *model.Payment,) error {

	return r.db.WithContext(ctx).
		Model(&model.Payment{}).
		Where(
			"payment_id = ?",
			payment.PaymentID,
		).
		Updates(map[string]interface{}{
			"payment_type":       payment.PaymentType,
			"amount":             payment.Amount,
			"currency":           payment.Currency,
			"status":             payment.Status,
			"payment_method":     payment.PaymentMethod,
			"external_reference": payment.ExternalReference,
			"scheduled_at":       payment.ScheduledAt,
			"processed_at":       payment.ProcessedAt,
		}).
		Error
}

func (r *PaymentRepo) Delete(ctx context.Context,id int64,) error {

	return r.db.WithContext(ctx).
		Where("payment_id = ?", id).
		Delete(&model.Payment{}).
		Error
}