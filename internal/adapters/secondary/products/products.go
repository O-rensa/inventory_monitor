package as_products

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
	c_brand "github.com/o-rensa/iv/internal/core/brands"
	c_inventories "github.com/o-rensa/iv/internal/core/inventory"
	c_productcategory "github.com/o-rensa/iv/internal/core/productcategories"
	c_product "github.com/o-rensa/iv/internal/core/products"
	pp_product "github.com/o-rensa/iv/internal/ports/primary/products"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var wg = new(sync.WaitGroup)

type ProductStore struct {
	db *gorm.DB
}

func NewProductStore(db *sql.DB) (*ProductStore, error) {
	ps := &ProductStore{}

	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err == nil {
		ps.db = gdb
	}

	return ps, err
}

// private methods
func (ps *ProductStore) getItemCount(wg *sync.WaitGroup, prodID uuid.UUID, itemCount *uint64, err *error) {
	defer wg.Done()
	var inv c_inventories.Inventory
	res := ps.db.First(&inv, "productID = ?", prodID.String())
	if res.Error == nil {
		*itemCount = inv.ItemCount
	}
	*err = res.Error
}

func (ps *ProductStore) getProductCategoryName(wg *sync.WaitGroup, pCID uuid.UUID, pCName *string, err *error) {
	defer wg.Done()
	var pc c_productcategory.ProductCategory
	res := ps.db.First(&pc, "id = ?", pCID.String())
	if res.Error == nil {
		*pCName = pc.ProductCategoryName
	}
	*err = res.Error
}

func (ps *ProductStore) getBrandName(wg *sync.WaitGroup, bID uuid.UUID, bName *string, err *error) {
	defer wg.Done()
	var b c_brand.Brand
	res := ps.db.First(&b, "id = ?", bID.String())
	if res.Error == nil {
		*bName = b.BrandName
	}
	*err = res.Error
}

// public methods
func (ps *ProductStore) CreateProduct(input c_product.Product) (pp_product.ProductDto, error) {

	prodID := input.ModelID.ID
	dto := pp_product.ProductDto{}
	inv := c_inventories.Inventory{
		ProductID: prodID,
		ItemCount: 0,
	}

	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return dto, err
	}

	// create product
	if err := tx.Create(&input).Error; err != nil {
		tx.Rollback()
		return dto, err
	}

	// create inventory
	if err := tx.Create(&inv).Error; err != nil {
		tx.Rollback()
		return dto, err
	}

	wg.Add(2)

	// get product category name
	var pCName string
	var pCNameErr error
	go ps.getProductCategoryName(wg, input.ProductCategoryID, &pCName, &pCNameErr)

	// get brand name
	var bName string
	var bNameErr error
	go ps.getBrandName(wg, input.BrandID, &bName, &bNameErr)

	wg.Wait()

	if pCNameErr != nil {
		tx.Rollback()
		return dto, pCNameErr
	}

	if bNameErr != nil {
		tx.Rollback()
		return dto, bNameErr
	}

	// insert values to dto
	dto.ProductID = prodID.String()
	dto.ProductCategoryID = input.ProductCategoryID.String()
	dto.ProductCategoryName = pCName
	dto.BrandID = input.ID.String()
	dto.BrandName = bName
	dto.ProductName = input.ProductName
	dto.ItemCount = "0"

	return dto, tx.Commit().Error
}

func (ps *ProductStore) GetAllProducts() ([]pp_product.ProductDto, error) {
	pDtos := []pp_product.ProductDto{}
	var prods []c_product.Product
	res := ps.db.Find(&prods)
	var retErr error
	if res.Error == nil {
		retErr = res.Error
		for _, v := range prods {
			wg.Add(3)
			// get item count
			var iCount uint64
			var iCountErr error
			go ps.getItemCount(wg, v.ModelID.ID, &iCount, &iCountErr)

			// get product category name
			var pCName string
			var pCNameErr error
			go ps.getProductCategoryName(wg, v.ProductCategoryID, &pCName, &pCNameErr)

			// get brand name
			var bName string
			var bNameErr error
			go ps.getBrandName(wg, v.BrandID, &bName, &bNameErr)

			wg.Wait()

			if iCountErr != nil {
				retErr = iCountErr
				break
			}

			if pCNameErr != nil {
				retErr = pCNameErr
				break
			}

			if bNameErr != nil {
				retErr = bNameErr
				break
			}

			i := pp_product.ProductDto{
				ProductID:           v.ModelID.ID.String(),
				ProductCategoryID:   v.ProductCategoryID.String(),
				ProductCategoryName: pCName,
				BrandID:             v.BrandID.String(),
				BrandName:           bName,
				ProductName:         v.ProductName,
				ItemCount:           fmt.Sprint(iCount),
			}
			pDtos = append(pDtos, i)
		}
	}

	if retErr != nil {
		pDtos = []pp_product.ProductDto{}
	}

	return pDtos, retErr
}

func (ps *ProductStore) GetProductByID(iD uuid.UUID) (pp_product.ProductDto, error) {
	var dto pp_product.ProductDto
	var p c_product.Product
	res := ps.db.First(&p, "id = ?", iD.String())
	if res.Error == nil {
		wg.Add(3)
		// get item count
		var iCount uint64
		var iCountErr error
		go ps.getItemCount(wg, iD, &iCount, &iCountErr)

		// get product category name
		var pCName string
		var pCNameErr error
		go ps.getProductCategoryName(wg, p.ProductCategoryID, &pCName, &pCNameErr)

		// get brand name
		var bName string
		var bNameErr error
		go ps.getBrandName(wg, p.BrandID, &bName, &bNameErr)

		wg.Wait()

		if iCountErr != nil {
			return dto, iCountErr
		}

		if pCNameErr != nil {
			return dto, pCNameErr
		}

		if bNameErr != nil {
			return dto, bNameErr
		}

		dto = pp_product.ProductDto{
			ProductID:           p.ModelID.ID.String(),
			ProductCategoryID:   p.ProductCategoryID.String(),
			ProductCategoryName: pCName,
			BrandID:             p.BrandID.String(),
			BrandName:           bName,
			ProductName:         p.ProductName,
			ItemCount:           fmt.Sprint(iCount),
		}
	}

	return dto, res.Error
}

func (ps *ProductStore) UpdateProduct(input c_product.Product) (pp_product.ProductDto, error) {
	var dto pp_product.ProductDto
	var p c_product.Product

	res := ps.db.First(&p, "id = ?", input.ModelID.ID.String())
	if res.Error != nil {
		return dto, res.Error
	}

	updated := ps.db.Model(&p).Updates(map[string]interface{}{
		"productCategoryId": input.ProductCategoryID,
		"brandId":           input.BrandID,
		"productName":       input.ProductName,
	})
	if updated.Error == nil {
		wg.Add(3)
		// get item count
		var iCount uint64
		var iCountErr error
		go ps.getItemCount(wg, input.ModelID.ID, &iCount, &iCountErr)

		// get product category name
		var pCName string
		var pCNameErr error
		go ps.getProductCategoryName(wg, input.ProductCategoryID, &pCName, &pCNameErr)

		// get brand name
		var bName string
		var bNameErr error
		go ps.getBrandName(wg, input.BrandID, &bName, &bNameErr)

		wg.Wait()

		if iCountErr != nil {
			return dto, iCountErr
		}

		if pCNameErr != nil {
			return dto, pCNameErr
		}

		if bNameErr != nil {
			return dto, bNameErr
		}

		dto = pp_product.ProductDto{
			ProductID:           input.ModelID.ID.String(),
			ProductCategoryID:   input.ProductCategoryID.String(),
			ProductCategoryName: pCName,
			BrandID:             input.BrandID.String(),
			BrandName:           bName,
			ProductName:         input.ProductName,
			ItemCount:           fmt.Sprint(iCount),
		}
	}

	return dto, res.Error
}

func (ps *ProductStore) UpdateProductItemCount(prodID uuid.UUID, count uint64) error {
	item := ps.db.Model(&c_inventories.Inventory{}).Where("productID = ?", prodID).Update("itemCount", count)
	return item.Error
}

func (ps *ProductStore) DeleteProduct(iD uuid.UUID) error {
	tx := ps.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// delete product
	if err := tx.Delete(&c_product.Product{}, iD).Error; err != nil {
		tx.Rollback()
		return err
	}

	// delete inventory
	if err := tx.Delete(&c_inventories.Inventory{}, iD).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
