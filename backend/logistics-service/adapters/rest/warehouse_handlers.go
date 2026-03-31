package rest

import (
	"net/http"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

func NewCreateWarehouseHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestCreateWarehouse
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		warehouse, err := svc.CreateWarehouse(r.Context(), req)
		if err != nil {
			if err == coreErrors.ErrDuplicateWarehouse {
				sendError(log, w, http.StatusConflict, "warehouse with this office_from_id already exists")
				return
			}
			sendError(log, w, http.StatusInternalServerError, "failed to create warehouse")
			return
		}

		sendCreated(log, w, warehouse)
	}
}

func NewGetWarehousesHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := svc.GetWarehouses(r.Context())
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to get warehouses")
			return
		}
		sendOK(log, w, list)
	}
}

func NewDeleteWarehouseHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("warehouse_id")
		if id == "" {
			sendError(log, w, http.StatusBadRequest, "warehouse_id is required")
			return
		}

		if err := svc.DeleteWarehouse(r.Context(), id); err != nil {
			if err == coreErrors.ErrNotFoundWarehouse {
				sendError(log, w, http.StatusNotFound, "warehouse not found")
				return
			}
			sendError(log, w, http.StatusInternalServerError, "failed to delete warehouse")
			return
		}

		sendNoContent(w)
	}
}