package data

import (
	"payments/internal/conf"
    "gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewPaymentRepo,
)

type Data struct {
	DB *gorm.DB
}

func NewData(c *conf.Data) (*Data, func(), error) {
	Db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{
		DB: Db,
	}, cleanup, nil
}
