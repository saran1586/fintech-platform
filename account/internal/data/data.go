package data

import (
	"account/internal/conf"
	"log/slog"

	"github.com/google/wire"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewDB,
	NewAccountRepo,
)

type Data struct {
	DB *gorm.DB
}

func NewData(c *conf.Data) (*Data, func(), error) {
	db, err := gorm.Open(postgres.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		slog.Info("closing database")
	}

	return &Data{
		DB: db,
	}, cleanup, nil
}

func NewDB(data *Data) *gorm.DB {
	return data.DB
}