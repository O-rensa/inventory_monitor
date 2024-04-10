package as_brand

import (
	"database/sql"

	"github.com/google/uuid"
	c_brand "github.com/o-rensa/iv/internal/core/brands"
	pp_brand "github.com/o-rensa/iv/internal/ports/primary/brands"
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

func (bs *BrandStore) CreateBrand(brand c_brand.Brand) (pp_brand.BrandDto, error) {
	res := bs.db.Create(brand)
	var dto pp_brand.BrandDto
	if res.Error == nil {
		dto = pp_brand.BrandDto{
			BrandID:   brand.ModelID.ID.String(),
			BrandName: brand.BrandName,
		}
	}

	return dto, res.Error
}

func (bs *BrandStore) GetAllBrands() ([]pp_brand.BrandDto, error) {
	var brands []c_brand.Brand
	dto := []pp_brand.BrandDto{}
	res := bs.db.Find(&brands)
	if res.Error == nil {
		for _, v := range brands {
			i := pp_brand.BrandDto{
				BrandID:   v.ModelID.ID.String(),
				BrandName: v.BrandName,
			}
			dto = append(dto, i)
		}
	}

	return dto, res.Error
}

func (bs *BrandStore) GetBrandByID(iD uuid.UUID) (pp_brand.BrandDto, error) {
	var dto pp_brand.BrandDto
	var brand c_brand.Brand
	res := bs.db.First(&brand, "id = ?", iD.String())
	if res.Error == nil {
		dto = pp_brand.BrandDto{
			BrandID:   brand.ModelID.ID.String(),
			BrandName: brand.BrandName,
		}
	}
	return dto, res.Error
}

func (bs *BrandStore) UpdateBrand(input c_brand.Brand) (pp_brand.BrandDto, error) {
	var dto pp_brand.BrandDto
	var brand c_brand.Brand

	// get brand by id
	gdb := bs.db.First(&brand, "id = ?", input.ModelID.ID.String())
	if gdb.Error != nil {
		return dto, gdb.Error
	}

	// update brand item
	res := bs.db.Model(&brand).Updates(map[string]interface{}{
		"brandName": input.BrandName,
	})
	if res.Error == nil {
		dto = pp_brand.BrandDto{
			BrandID:   input.ID.String(),
			BrandName: input.BrandName,
		}
	}
	return dto, res.Error
}

func (bs *BrandStore) DeleteBrand(iD uuid.UUID) error {
	return (bs.db.Delete(&c_brand.Brand{}, iD)).Error
}
