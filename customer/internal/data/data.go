package data

import (
	"customer/internal/conf"
    "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/go-kratos/kratos/v3/log"
	"github.com/google/wire"

)


var ProviderSet = wire.NewSet(
	NewData,
	NewCustomerRepo,
)
type Data struct {
	DB *gorm.DB
}


// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
db,err := gorm.Open(postgres.Open(c.Database.Source),&gorm.Config{})
if err != nil{
	panic(err)
}

	cleanup := func() {
		log.Info("closing the data resources")
	}
	return &Data{
DB:db,
	}, cleanup, nil
}
