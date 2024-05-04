package c_deliveries

import (
	"time"

	"github.com/google/uuid"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
)

type DeliveryDetails struct {
	DeliveryID uuid.UUID `gorm:"PrimaryKey;column:deliveryID"`
	*c_sharedTypes.BaseModel
	DeliveryDate   time.Time      `gorm:"column:deliveryDate"`
	ItemsDelivered []DeliveryItem `gorm:"foreignKey:DeliveryID;references:DeliveryID"`
}

type DeliveryItem struct {
	DeliveryID uuid.UUID `gorm:"PrimaryKey;column:deliveryID"`
	ProductID  uuid.UUID `gorm:"column:productID"`
	ItemCount  uint64    `gorm:"column:itemCount"`
}
