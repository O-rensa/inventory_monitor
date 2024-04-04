package c_products

import (
	"github.com/google/uuid"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
)

type Product struct {
	*c_sharedTypes.ModelID
	*c_sharedTypes.BaseModel
	ProductCategoryId uuid.UUID `gorm:"productCategoryId"`
	BrandId           uuid.UUID `gorm:"brandId"`
	ProductName       string    `gorm:"productName"`
}
