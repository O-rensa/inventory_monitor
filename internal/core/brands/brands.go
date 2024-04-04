package c_brands

import c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"

type Brands struct {
	*c_sharedTypes.ModelID
	*c_sharedTypes.BaseModel
	BrandName string `gorm:"brandName"`
}
