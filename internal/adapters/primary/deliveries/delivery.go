package ap_delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	ap_auth "github.com/o-rensa/iv/internal/adapters/primary/auth"
	ap_shared "github.com/o-rensa/iv/internal/adapters/primary/shared"
	c_deliveries "github.com/o-rensa/iv/internal/core/deliveries"
	pp_delivery "github.com/o-rensa/iv/internal/ports/primary/deliveries"
	ps_admin "github.com/o-rensa/iv/internal/ports/secondary/admin"
	ps_delivery "github.com/o-rensa/iv/internal/ports/secondary/deliveries"
)

type DeliveryHandler struct {
	deliveryStore ps_delivery.DeliveryStore
	adminStore    ps_admin.AdminStore
}

func NewDeliveryHandler(deliveryStore ps_delivery.DeliveryStore, adminStore ps_admin.AdminStore) *DeliveryHandler {
	dh := &DeliveryHandler{
		deliveryStore: deliveryStore,
		adminStore:    adminStore,
	}

	return dh
}

func (dh *DeliveryHandler) RegisterRoutes(router *mux.Router) {
	sr := router.PathPrefix("/delivery").Subrouter()

	sr.HandleFunc("/create", ap_auth.WithJWTAuth(dh.CreateDeliveryTransaction, dh.adminStore)).Methods(http.MethodPost)
	sr.HandleFunc("/delete/{iD}", ap_auth.WithJWTAuth(dh.DeleteDeliveryTransaction, dh.adminStore)).Methods(http.MethodDelete)
}

func (dh *DeliveryHandler) CreateDeliveryTransaction(w http.ResponseWriter, r *http.Request) {
	var payload pp_delivery.CreateDeliveryPayload
	if err := ap_shared.ParseJSON(r, &payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := ap_shared.Validate.Struct(payload); err != nil {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", err))
		return
	}

	delId := uuid.New()

	d := c_deliveries.DeliveryDetails{
		DeliveryID:   delId,
		DeliveryDate: payload.DeliveryDate,
	}

	for _, item := range payload.ItemsDelivered {
		ic, _ := strconv.ParseUint(item.ItemCount, 10, 64)
		i := c_deliveries.DeliveryItem{
			DeliveryID: delId,
			ProductID:  uuid.MustParse(item.ProductId),
			ItemCount:  ic,
		}
		d.ItemsDelivered = append(d.ItemsDelivered, i)
	}

	if err := dh.deliveryStore.CreateDeliveryTransaction(d); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusCreated, nil)
}

func (dh *DeliveryHandler) DeleteDeliveryTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	str, ok := vars["iD"]
	if !ok {
		ap_shared.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing delivery transaction id"))
		return
	}

	dID := uuid.MustParse(str)
	if err := dh.deliveryStore.DeleteDeliveryTransaction(dID); err != nil {
		ap_shared.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	ap_shared.WriteJSON(w, http.StatusOK, nil)
}
