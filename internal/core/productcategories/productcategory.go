package c_productcategory

import c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"

type ProductCategory struct {
	*c_sharedTypes.ModelID
	*c_sharedTypes.BaseModel
	ProductCategoryName string `gorm:"column:productCategoryName"`
}
