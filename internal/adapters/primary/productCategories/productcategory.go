package ap_productCategory

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	sr.HandleFunc("/getall", ap_auth.WithJWTAuth(pch.GetAllProductCategoryHandler, pch.adminStore)).Methods(http.MethodGet)
	sr.HandleFunc("/get/{iD}", ap_auth.WithJWTAuth(pch.GetProductCategoryByIDHandler, pch.adminStore)).Methods(http.MethodGet)
	sr.HandleFunc("/update", ap_auth.WithJWTAuth(pch.UpdateProductCategoryHandler, pch.adminStore)).Methods(http.MethodPut)
	sr.HandleFunc("/delete/{iD}", ap_auth.WithJWTAuth(pch.DeleteProductCategoryHandler, pch.adminStore)).Methods(http.MethodDelete)
}

func (pch *ProductCategoryHandler) CreateProductCategoryHandler(w http.ResponseWriter, r *http.Request) {
	// decode request and store into paylaod variable
	var payload pp_productCategory.CreateProductCategoryPayload
	if err := ap_shared.ParseJSON(r, &payload); err != nil {
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

func (pch *ProductCategoryHandler) GetAllProductCategoryHandler(w http.ResponseWriter, r *http.Request) {
	res, err := pch.productCategoryStore.GetAllProductCategory()
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (pch *ProductCategoryHandler) GetProductCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product category id"))
		return
	}

	pCID := uuid.MustParse(str)
	res, err := pch.productCategoryStore.GetProductCategoryByID(pCID)
	if err != nil {
		ap_shared.WriteError(w, http.StatusNotFound, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (pch *ProductCategoryHandler) UpdateProductCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var pl pp_productCategory.UpdateProductCategoryPayload
	if err := ap_shared.ParseJSON(r, &pl); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := ap_shared.Validate.Struct(pl); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	pCID := uuid.MustParse(pl.ProductCategoryID)
	pc := c_productcategory.ProductCategory{
		ModelID:             &c_sharedTypes.ModelID{ID: pCID},
		ProductCategoryName: pl.ProductCategoryName,
	}

	res, err := pch.productCategoryStore.UpdateProductCategory(pc)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (pch *ProductCategoryHandler) DeleteProductCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product category id"))
		return
	}

	pCID := uuid.MustParse(str)
	if err := pch.productCategoryStore.DeleteProductCategory(pCID); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, nil)
}
