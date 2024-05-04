package c_sales

import (
	"time"

	"github.com/google/uuid"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
)

type SaleDetails struct {
	SaleID uuid.UUID `gorm:"PrimaryKey;column:saleID"`
	*c_sharedTypes.BaseModel
	SaleDate  time.Time  `gorm:"column:saleDate"`
	ItemsSold []SaleItem `gorm:"foreignKey:SaleID;references:SaleID"`
}

type SaleItem struct {
	SaleID    uuid.UUID `gorm:"PrimaryKey;column:saleID"`
	ProductId uuid.UUID `gorm:"column:productID"`
	ItemCount uint64    `gorm:"column:itemCount"`
}
