package ap_product

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	c_product "github.com/o-rensa/iv/internal/core/products"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
	pp_product "github.com/o-rensa/iv/internal/ports/primary/products"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	ps_product "github.com/o-rensa/iv/internal/ports/secondary/products"
)

type ProductHandler struct {
	productStore ps_product.ProductStore
	adminStore   ps_admin.AdminStore
}

func NewProductHandler(productStore ps_product.ProductStore, adminStore ps_admin.AdminStore) *ProductHandler {
	ph := &ProductHandler{
		productStore: productStore,
		adminStore:   adminStore,
	}
	return ph
}

func (ph *ProductHandler) RegisterRoutes(router *mux.Router) {
	sr := router.PathPrefix("/products").Subrouter()

	sr.HandleFunc("/create", ap_auth.WithJWTAuth(ph.CreateProductHandler, ph.adminStore)).Methods(http.MethodPost)
	sr.HandleFunc("/getall", ap_auth.WithJWTAuth(ph.GetAllProductHandler, ph.adminStore)).Methods(http.MethodGet)
	sr.HandleFunc("/get/{iD}", ap_auth.WithJWTAuth(ph.GetProductByIDHandler, ph.adminStore)).Methods(http.MethodGet)
	sr.HandleFunc("/update", ap_auth.WithJWTAuth(ph.UpdateProductHandler, ph.adminStore)).Methods(http.MethodPut)
	sr.HandleFunc("/delete/{iD}", ap_auth.WithJWTAuth(ph.DeleteProductHandler, ph.adminStore)).Methods(http.MethodDelete)
}

func (ph *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var payload pp_product.CreateProductPayload
	if err := ap_shared.ParseJSON(r, &payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := ap_shared.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	p := c_product.Product{
		ModelID:           c_sharedTypes.NewModelID(),
		ProductCategoryID: uuid.MustParse(payload.ProductCategoryID),
		BrandID:           uuid.MustParse(payload.BrandID),
		ProductName:       payload.ProductName,
	}

	dto, err := ph.productStore.CreateProduct(p)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusCreated, dto)
}

func (ph *ProductHandler) GetAllProductHandler(w http.ResponseWriter, r *http.Request) {
	res, err := ph.productStore.GetAllProducts()
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (ph *ProductHandler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	str, ok := v["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product id"))
		return
	}

	pID := uuid.MustParse(str)
	res, err := ph.productStore.GetProductByID(pID)
	if err != nil {
		ap_shared.WriteError(w, http.StatusNotFound, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (ph *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	var pl pp_product.UpdateProductPayload
	if err := ap_shared.ParseJSON(r, &pl); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := ap_shared.Validate.Struct(pl); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	pID := uuid.MustParse(pl.ProductID)
	p := c_product.Product{
		ModelID:           &c_sharedTypes.ModelID{ID: pID},
		ProductCategoryID: uuid.MustParse(pl.ProductCategoryID),
		BrandID:           uuid.MustParse(pl.BrandID),
		ProductName:       pl.ProductName,
	}

	res, err := ph.productStore.UpdateProduct(p)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (ph *ProductHandler) UpdateProductItemCountHandler(w http.ResponseWriter, r *http.Request) {
	var pl pp_product.UpdateProductItemCount
	if err := ap_shared.ParseJSON(r, &pl); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := ap_shared.Validate.Struct(pl); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	pID := uuid.MustParse(pl.ProductID)
	iC, err := strconv.ParseUint(pl.ItemCount, 10, 64)
	if err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	retErr := ph.productStore.UpdateProductItemCount(pID, iC)
	if retErr != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, retErr)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, nil)
}

func (ph *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	v := mux.Vars(r)
	str, ok := v["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing product id"))
		return
	}

	pID := uuid.MustParse(str)
	if err := ph.productStore.DeleteProduct(pID); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, nil)
}
