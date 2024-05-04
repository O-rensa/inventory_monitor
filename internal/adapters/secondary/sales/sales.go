package as_sales

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func (ss *SalesStore) CreateNewSalesTransaction() {

}

func (ss *SalesStore) DeleteSalesTransaction() {

}
