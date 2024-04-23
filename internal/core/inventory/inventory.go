package c_inventories

import (
	"github.com/google/uuid"
	c_products "github.com/o-rensa/iv/internal/core/products"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
)

type Inventory struct {
	*c_sharedTypes.BaseModel
	ProductID uuid.UUID          `gorm:"primaryKey;column:productID"`
	Product   c_products.Product `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ItemCount uint               `gorm:"column:itemCount"`
}
