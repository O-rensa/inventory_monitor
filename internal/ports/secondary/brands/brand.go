package ps_brand

import (
	"github.com/google/uuid"
	c_brand "github.com/o-rensa/iv/internal/core/brands"
	pp_brand "github.com/o-rensa/iv/internal/ports/primary/brands"
)

type BrandStore interface {
	CreateBrand(c_brand.Brand) (pp_brand.BrandDto, error)
	GetAllBrands() ([]pp_brand.BrandDto, error)
	GetBrandByID(uuid.UUID) (pp_brand.BrandDto, error)
	UpdateBrand(c_brand.Brand) (pp_brand.BrandDto, error)
	DeleteBrand(uuid.UUID) error
}
