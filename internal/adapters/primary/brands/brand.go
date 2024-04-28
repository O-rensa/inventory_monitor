package ap_brand

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	c_brand "github.com/o-rensa/iv/internal/core/brands"
	c_sharedTypes "github.com/o-rensa/iv/internal/core/sharedtypes"
	pp_brand "github.com/o-rensa/iv/internal/ports/primary/brands"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	ps_brand "github.com/o-rensa/iv/internal/ports/secondary/brands"
)

type BrandHandler struct {
	brandStore ps_brand.BrandStore
	adminStore ps_admin.AdminStore
}

func NewBrandHandler(brandStore ps_brand.BrandStore, adminStore ps_admin.AdminStore) *BrandHandler {
	bh := &BrandHandler{brandStore: brandStore, adminStore: adminStore}
	return bh
}

func (bh *BrandHandler) RegisterRoutes(router *mux.Router) {
	sr := router.PathPrefix("/brands").Subrouter()

	sr.HandleFunc("/create", ap_auth.WithJWTAuth(bh.CreateBrandHandler, bh.adminStore)).Methods(http.MethodPost)
	sr.HandleFunc("/getall", ap_auth.WithJWTAuth(bh.GetAllBrandsHandler, bh.adminStore)).Methods(http.MethodGet)
	sr.HandleFunc("/get/{iD}", ap_auth.WithJWTAuth(bh.GetBrandByIDHandler, bh.adminStore)).Methods(http.MethodGet)
	sr.HandleFunc("/update", ap_auth.WithJWTAuth(bh.UpdateBrandHandler, bh.adminStore)).Methods(http.MethodPut)
	sr.HandleFunc("/delete/{iD}", ap_auth.WithJWTAuth(bh.DeleteBrandHandler, bh.adminStore)).Methods(http.MethodDelete)
}

func (bh *BrandHandler) CreateBrandHandler(w http.ResponseWriter, r *http.Request) {
	// decode request and store into payload variable
	var payload pp_brand.CreateBrandPayload
	if err := ap_shared.ParseJSON(r, &payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// validate payload
	if err := ap_shared.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
	}

	// create brand
	brand := c_brand.Brand{
		ModelID:   c_sharedTypes.NewModelID(),
		BrandName: payload.BrandName,
	}
	dto, err := bh.brandStore.CreateBrand(brand)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
	}

	ap_shared.WriteJSON(w, http.StatusCreated, dto)
}

func (bh *BrandHandler) GetAllBrandsHandler(w http.ResponseWriter, r *http.Request) {
	res, err := bh.brandStore.GetAllBrands()
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (bh *BrandHandler) GetBrandByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing brand id"))
		return
	}

	// convert id into uuid
	bID := uuid.MustParse(str)
	res, err := bh.brandStore.GetBrandByID(bID)
	if err != nil {
		ap_shared.WriteError(w, http.StatusNotFound, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (bh *BrandHandler) UpdateBrandHandler(w http.ResponseWriter, r *http.Request) {
	var pl pp_brand.UpdateBrandPayload
	if err := ap_shared.ParseJSON(r, &pl); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := ap_shared.Validate.Struct(pl); err != nil {
		errors := err.(validator.ValidationErrors)
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	bID := uuid.MustParse(pl.BrandID)
	b := c_brand.Brand{
		ModelID:   &c_sharedTypes.ModelID{ID: bID},
		BrandName: pl.BrandName,
	}

	res, err := bh.brandStore.UpdateBrand(b)
	if err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
	}

	ap_shared.WriteJSON(w, http.StatusOK, res)
}

func (bh *BrandHandler) DeleteBrandHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing brand id"))
		return
	}

	bID := uuid.MustParse(str)
	if err := bh.brandStore.DeleteBrand(bID); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, fmt.Errorf("something went wrong: %v", err))
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, nil)
}
