package main

import (
	"log"

	c_admin "github.com/o-rensa/iv/internal/core/admin"
	c_brands "github.com/o-rensa/iv/internal/core/brands"
	c_productcategories "github.com/o-rensa/iv/internal/core/productcategories"
	c_products "github.com/o-rensa/iv/internal/core/products"
	"github.com/o-rensa/iv/pkg/initializers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var modelsToMigrate = []interface{}{
	&c_admin.Admin{},
	&c_brands.Brand{},
	&c_productcategories.ProductCategory{},
	&c_products.Product{},
}

var db *gorm.DB

func init() {
	initializers.LoadDotEnv()
	sqlDB, err := initializers.InitializePostGres()
	if err != nil {
		log.Fatal(err)
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db.AutoMigrate(modelsToMigrate...)
}
