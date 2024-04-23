package as_products

import (
	"database/sql"

	"github.com/google/uuid"
	c_inventories "github.com/o-rensa/iv/internal/core/inventory"
	c_products "github.com/o-rensa/iv/internal/core/products"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductStore struct {
	db *gorm.DB
}

func NewProductStore(db *sql.DB) (*ProductStore, error) {
	ps := &ProductStore{}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		ps.db = gdb
	}

	return ps, err
}

func (ps *ProductStore) CreateProduct(input c_products.Product) error {

	prodID := input.ModelID.ID

	inv := &c_inventories.Inventory{
		ProductID: prodID,
		ItemCount: 0,
	}

	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// create product
	if err := tx.Create(input).Error; err != nil {
		tx.Rollback()
		return err
	}

	// create inventory
	if err := tx.Create(inv).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ps *ProductStore) DeleteProduct(iD uuid.UUID) error {
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// delete product
	if err := tx.Delete(&c_products.Product{}, iD).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete inventory
	if err := tx.Delete(&c_inventories.Inventory{}, iD).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
