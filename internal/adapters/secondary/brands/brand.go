package ps_brand

import (
	"database/sql"

	"github.com/google/uuid"
	c_brands "github.com/o-rensa/iv/internal/core/brands"
	pp_brands "github.com/o-rensa/iv/internal/ports/primary/brands"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BrandStore struct {
	db *gorm.DB
}

func NewBrandStore(db *sql.DB) (*BrandStore, error) {
	bs := &BrandStore{}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		bs.db = gormDB
	}

	return bs, err
}

func (bs *BrandStore) CreateBrand(brand *c_brands.Brand) (*pp_brands.BrandDto, error) {
	res := bs.db.Create(brand)
	var dto *pp_brands.BrandDto
	if res.Error == nil {
		dto = &pp_brands.BrandDto{
			BrandID:   brand.ModelID.ID.String(),
			BrandName: brand.BrandName,
		}
	}

	return dto, res.Error
}

func (bs *BrandStore) GetAllBrands() ([]pp_brands.BrandDto, error) {
	var brands []c_brands.Brand
	dto := []pp_brands.BrandDto{}
	res := bs.db.Find(&brands)
	if res.Error == nil {
		for _, v := range brands {
			i := pp_brands.BrandDto{
				BrandID:   v.ModelID.ID.String(),
				BrandName: v.BrandName,
			}
			dto = append(dto, i)
		}
	}

	return dto, res.Error
}

func (bs *BrandStore) GetBrandByID(iD uuid.UUID) (*pp_brands.BrandDto, error) {
	var dto *pp_brands.BrandDto
	var brand c_brands.Brand
	res := bs.db.First(&brand, "id = ?", iD.String())
	if res.Error == nil {
		dto = &pp_brands.BrandDto{
			BrandID:   brand.ModelID.ID.String(),
			BrandName: brand.BrandName,
		}
	}
	return dto, res.Error
}

func (bs *BrandStore) UpdateBrand(brand *c_brands.Brand) (*pp_brands.BrandDto, error) {
	var dto *pp_brands.BrandDto
	res := bs.db.Save(brand)
	if res.Error == nil {
		dto = &pp_brands.BrandDto{
			BrandID:   brand.ID.String(),
			BrandName: brand.BrandName,
		}
	}
	return dto, res.Error
}

func (bs *BrandStore) DeleteBrandByID(iD uuid.UUID) error {
	return (bs.db.Delete(&c_brands.Brand{}, iD)).Error
}
