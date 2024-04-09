package c_brands

import c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"

type Brand struct {
	*c_sharedTypes.ModelID
	*c_sharedTypes.BaseModel
	BrandName string `gorm:"column:brandName"`
}
