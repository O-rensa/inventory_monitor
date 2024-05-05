package as_sales

import (
	"database/sql"

	"github.com/google/uuid"
	c_sales "github.com/o-rensa/iv/internal/core/sales"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SalesStore struct {
	db *gorm.DB
}

func NewSalesStore(db *sql.DB) (*SalesStore, error) {
	ss := &SalesStore{}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		ss.db = gdb
	}

	return ss, err
}

func (ss *SalesStore) CreateNewSalesTransaction(input c_sales.SaleDetails) error {
	tx := ss.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// create sales details
	if err := tx.Omit(clause.Associations).Create(&input).Error; err != nil {
		tx.Rollback()
		return err
	}

	// create sale items
	for _, item := range input.ItemsSold {
		si := &c_sales.SaleItem{
			SaleID:    item.SaleID,
			ProductId: item.ProductId,
			ItemCount: item.ItemCount,
		}

		if err := tx.Create(si).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (ss *SalesStore) DeleteSalesTransaction(iD uuid.UUID) error {
	tx := ss.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// delete sale detail
	if err := tx.Where("saleID = ?", iD).Delete(&c_sales.SaleDetails{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete sale items
	if err := tx.Where("saleID = ?", iD).Delete(&c_sales.SaleItem{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
