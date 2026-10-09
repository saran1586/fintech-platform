package data

import (
	"context"

	"customer/internal/biz"
	"customer/internal/model"

	"gorm.io/gorm"
)

type CustomerRepo struct {
	db *gorm.DB
}

func NewCustomerRepo(data *Data) biz.CustomerRepo {
	return &CustomerRepo{
		db: data.DB,
	}
}

func (r *CustomerRepo) Create(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Create(customer).Error
}

func (r *CustomerRepo) GetByID(ctx context.Context, id int64) (*model.Customer, error) {
	var customer model.Customer

	err := r.db.WithContext(ctx).
		Where("customer_id = ?", id).
		First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *CustomerRepo) GetAll(ctx context.Context) ([]*model.Customer, error) {
	var customers []*model.Customer

	err := r.db.WithContext(ctx).
		Find(&customers).Error

	if err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerRepo) Update(ctx context.Context, customer *model.Customer) error {
	return r.db.WithContext(ctx).Save(customer).Error
}