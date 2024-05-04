package as_deliveries

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DeliveryStore struct {
	db *gorm.DB
}

func NewDeliveryStore(db *sql.DB) (*DeliveryStore, error) {
	ds := &DeliveryStore{}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		ds.db = gdb
	}

	return ds, err
}

func (ds *DeliveryStore) CreateDeliveryTransaction() {

}

func (ds *DeliveryStore) DeleteDeliveryTransaction() {

}
