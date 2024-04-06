package ps_brands

import (
	"github.com/google/uuid"
	c_brands "github.com/o-rensa/iv/internal/core/brands"
	pp_brands "github.com/o-rensa/iv/internal/ports/primary/brands"
)

type BrandStore interface {
	CreateBrand(*c_brands.Brand) (*pp_brands.BrandDto, error)
	GetAllBrands() ([]pp_brands.BrandDto, error)
	GetBrandByID(uuid.UUID) (*pp_brands.BrandDto, error)
	UpdateBrand(*c_brands.Brand) (*pp_brands.BrandDto, error)
	DeleteBrand(uuid.UUID) error
}
