package data

import (
	"context"

	"account/internal/biz"
	"account/internal/model"

	"gorm.io/gorm"
)

type AccountRepo struct {
	db *gorm.DB
}

func NewAccountRepo(db *gorm.DB) biz.AccountRepo {
	return &AccountRepo{
		db: db,
	}
}

func (r *AccountRepo) Create(ctx context.Context, account *model.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

 func (r *AccountRepo) GetByID(ctx context.Context, id int64) (*model.Account, error) {
	var account model.Account

	err := r.db.WithContext(ctx).
		Where("account_id = ?", id).
		First(&account).Error

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *AccountRepo) GetAll(ctx context.Context) ([]*model.Account, error) {
	var accounts []*model.Account

	err := r.db.WithContext(ctx).Find(&accounts).Error
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (r *AccountRepo) Update(ctx context.Context,account *model.Account,) error {
	return r.db.WithContext(ctx).
		Save(account).Error
}
func (r *AccountRepo) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).
		Model(&model.Account{}).
		Where("account_id = ?", id).
		Update("status", "closed").Error
}