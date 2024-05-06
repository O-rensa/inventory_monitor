package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	ap_admin "github.com/o-rensa/iv/internal/adapters/primary/admin"
	ap_brand "github.com/o-rensa/iv/internal/adapters/primary/brands"
	ap_delivery "github.com/o-rensa/iv/internal/adapters/primary/deliveries"
	ap_product "github.com/o-rensa/iv/internal/adapters/primary/product"
	ap_productCategory "github.com/o-rensa/iv/internal/adapters/primary/productCategories"
	ap_sales "github.com/o-rensa/iv/internal/adapters/primary/sales"
	as_admin "github.com/o-rensa/iv/internal/adapters/secondary/admin"
	as_brand "github.com/o-rensa/iv/internal/adapters/secondary/brands"
	as_deliveries "github.com/o-rensa/iv/internal/adapters/secondary/deliveries"
	as_productCategory "github.com/o-rensa/iv/internal/adapters/secondary/productcategories"
	as_products "github.com/o-rensa/iv/internal/adapters/secondary/products"
	as_sales "github.com/o-rensa/iv/internal/adapters/secondary/sales"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: "localhost:" + addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// admin
	adminStore, _ := as_admin.NewAdminStore(s.db)
	adminHandler := ap_admin.NewAdminHandler(adminStore)
	adminHandler.RegisterRoutes(subrouter)

	// brand
	brandStore, _ := as_brand.NewBrandStore(s.db)
	brandhandler := ap_brand.NewBrandHandler(brandStore, adminStore)
	brandhandler.RegisterRoutes(subrouter)

	// product category
	productCategoryStore, _ := as_productCategory.NewProductCategoryStore(s.db)
	productCategoryHandler := ap_productCategory.NewProductCategoryHandler(productCategoryStore, adminStore)
	productCategoryHandler.RegisterRoutes(subrouter)

	// product
	productStore, _ := as_products.NewProductStore(s.db)
	productHandler := ap_product.NewProductHandler(productStore, adminStore)
	productHandler.RegisterRoutes(subrouter)

	// sales transactions
	salesStore, _ := as_sales.NewSalesStore(s.db)
	salesHandler := ap_sales.NewSalesHandler(salesStore, adminStore)
	salesHandler.RegisterRoutes(subrouter)

	// delivery transactions
	deliveryStore, _ := as_deliveries.NewDeliveryStore(s.db)
	deliveryHandler := ap_delivery.NewDeliveryHandler(deliveryStore, adminStore)
	deliveryHandler.RegisterRoutes(subrouter)

	log.Println("Listening on", s.addr)
	return http.ListenAndServe(s.addr, router)
}
