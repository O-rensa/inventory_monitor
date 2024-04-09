package as_productCategories

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ProductCategoryStore struct {
	db *gorm.DB
}

func NewProductCategoryStore(db *sql.DB) (*ProductCategoryStore, error) {
	pcs := &ProductCategoryStore{}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		pcs.db = gdb
	}

	return pcs, err
}
