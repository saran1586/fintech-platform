package data

import (
	"context"

	"card/internal/biz"
	"card/internal/model"

	"gorm.io/gorm"
)

type CardRepo struct {
	db *gorm.DB
}

func NewCardRepo(db *gorm.DB) biz.CardRepo {
	return &CardRepo{
		db: db,
	}
}

func (r *CardRepo) Create(ctx context.Context, card *model.Card) error {
	return r.db.WithContext(ctx).Create(card).Error
}

func (r *CardRepo) GetByID(ctx context.Context, id int64) (*model.Card, error) {
	var card model.Card

	err := r.db.WithContext(ctx).
		Where("card_id = ?", id).
		First(&card).Error

	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (r *CardRepo) GetAll(ctx context.Context,) ([]*model.Card, error) {

	var cards []*model.Card

	err := r.db.WithContext(ctx).
		Find(&cards).Error

	if err != nil {
		return nil, err
	}

	return cards, nil
}


func (r *CardRepo) Update(ctx context.Context, card *model.Card) error {
	return r.db.WithContext(ctx).
		Save(card).Error
}

func (r *CardRepo) Delete(ctx context.Context,id int64,) error {
	return r.db.WithContext(ctx).
		Delete(&model.Card{}, id).Error
}