package ap_productCategories

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	pp_productCategories "github.com/o-rensa/iv/internal/ports/primary/productcategories"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	ps_productCategories "github.com/o-rensa/iv/internal/ports/secondary/productcategories"
)

type ProductCategoryHandler struct {
	productCategoryStore ps_productCategories.ProductCategoryStore
	adminStore           ps_admin.AdminStore
}

func NewProductCategoryHandler(productCategoryStore ps_productCategories.ProductCategoryStore, adminStore ps_admin.AdminStore) *ProductCategoryHandler {
	pch := &ProductCategoryHandler{
		productCategoryStore: productCategoryStore,
		adminStore:           adminStore,
	}
	return pch
}

func (pch *ProductCategoryHandler) RegisterRoutes(router *mux.Router) {
	sr := router.PathPrefix("/productcategories").Subrouter()

	sr.HandleFunc("/create", ap_auth.WithJWTAuth(pch.CreateProductCategoryHandler, pch.adminStore)).Methods(http.MethodPost)
}

func (pch *ProductCategoryHandler) CreateProductCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// decode request and store into paylaod variable
	var payload pp_productCategories.CreateProductCategoryPayload
	if err := ap_shared.ParseJSON(r, &payload); err == nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := ap_shared.Validate.Struct(payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	ap_shared.WriteJSON(w, http.StatusCreated, nil)
}
