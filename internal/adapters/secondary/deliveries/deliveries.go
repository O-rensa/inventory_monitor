package as_deliveries

import (
	"database/sql"

	"github.com/google/uuid"
	c_deliveries "github.com/o-rensa/iv/internal/core/deliveries"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (ds *DeliveryStore) CreateDeliveryTransaction(input c_deliveries.DeliveryDetails) error {
	tx := ds.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// create delivery details
	if err := tx.Omit(clause.Associations).Create(&input).Error; err != nil {
		tx.Rollback()
		return err
	}

	// create delivery items
	for _, item := range input.ItemsDelivered {
		di := &c_deliveries.DeliveryItem{
			DeliveryID: item.DeliveryID,
			ProductID:  item.ProductID,
			ItemCount:  item.ItemCount,
		}

		if err := tx.Create(di).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (ds *DeliveryStore) DeleteDeliveryTransaction(iD uuid.UUID) error {
	tx := ds.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// delete delivery detail
	if err := tx.Where("deliveryID = ?", iD).Delete(&c_deliveries.DeliveryDetails{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete delivery items
	if err := tx.Where("deliveryID = ?", iD).Delete(&c_deliveries.DeliveryItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
