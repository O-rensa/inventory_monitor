package c_product

import (
	"github.com/google/uuid"
	c_brands "github.com/o-rensa/iv/internal/core/brands"
	c_productcategories "github.com/o-rensa/iv/internal/core/productcategories"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
)

type Product struct {
	*c_sharedTypes.ModelID
	*c_sharedTypes.BaseModel
	//product category
	ProductCategoryID uuid.UUID                           `gorm:"column:productCategoryId"`
	ProductCategory   c_productcategories.ProductCategory `gorm:"foreignKey:ProductCategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	// brand
	BrandID uuid.UUID      `gorm:"column:brandId"`
	Brand   c_brands.Brand `gorm:"foreignKey:BrandID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ProductName string `gorm:"column:productName"`
}
