package ap_productCategory

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	c_productcategory "github.com/o-rensa/iv/internal/core/productcategories"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
	pp_productCategory "github.com/o-rensa/iv/internal/ports/primary/productcategories"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	ps_productCategory "github.com/o-rensa/iv/internal/ports/secondary/productcategories"
)

type ProductCategoryHandler struct {
	productCategoryStore ps_productCategory.ProductCategoryStore
	adminStore           ps_admin.AdminStore
}

func NewProductCategoryHandler(productCategoryStore ps_productCategory.ProductCategoryStore, adminStore ps_admin.AdminStore) *ProductCategoryHandler {
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
	var payload pp_productCategory.CreateProductCategoryPayload
	if err := ap_shared.ParseJSON(r, &payload); err == nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := ap_shared.Validate.Struct(payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	// create Product Category
	pc := c_productcategory.ProductCategory{
		ModelID:             c_sharedTypes.NewModelID(),
		ProductCategoryName: payload.ProductCategoryName,
	}
	dto, err := pch.productCategoryStore.CreateProductCategory(pc)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
	}

	ap_shared.WriteJSON(w, http.StatusCreated, dto)
}
