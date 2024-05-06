package ap_sales

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	c_sales "github.com/o-rensa/iv/internal/core/sales"
	pp_sales "github.com/o-rensa/iv/internal/ports/primary/sales"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	ps_sales "github.com/o-rensa/iv/internal/ports/secondary/sales"
)

type SalesHandler struct {
	salesStore ps_sales.SalesStore
	adminStore ps_admin.AdminStore
}

func NewSalesHandler(salesStore ps_sales.SalesStore, adminStore ps_admin.AdminStore) *SalesHandler {
	sh := &SalesHandler{
		salesStore: salesStore,
		adminStore: adminStore,
	}

	return sh
}

func (sh *SalesHandler) RegisterRoutes(router *mux.Router) {
	sr := router.PathPrefix("/sales").Subrouter()

	sr.HandleFunc("/create", ap_auth.WithJWTAuth(sh.CreateSalesTransaction, sh.adminStore)).Methods(http.MethodPost)
	sr.HandleFunc("/delete/{iD}", ap_auth.WithJWTAuth(sh.DeleteSalesTransaction, sh.adminStore)).Methods(http.MethodDelete)
}

func (sh *SalesHandler) CreateSalesTransaction(w http.ResponseWriter, r *http.Request) {
	var payload pp_sales.CreateSalesPayload
	if err := ap_shared.ParseJSON(r, &payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := ap_shared.Validate.Struct(payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	saleId := uuid.New()

	s := c_sales.SaleDetails{
		SaleID:   saleId,
		SaleDate: payload.SaleDate,
	}

	for _, item := range payload.ItemsSold {
		ic, _ := strconv.ParseUint(item.ItemCount, 10, 64)
		i := c_sales.SaleItem{
			SaleID:    saleId,
			ProductId: uuid.MustParse(item.ProductId),
			ItemCount: ic,
		}
		s.ItemsSold = append(s.ItemsSold, i)
	}

	if err := sh.salesStore.CreateNewSalesTransaction(s); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusCreated, nil)
}

func (sh *SalesHandler) DeleteSalesTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing sales transaction id"))
		return
	}

	sID := uuid.MustParse(str)
	if err := sh.salesStore.DeleteSalesTransaction(sID); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, nil)
}
