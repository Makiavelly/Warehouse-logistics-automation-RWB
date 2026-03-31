package rest

import (
	"net/http"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

func NewCreateRouteHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.PathValue("warehouse_id")
		if warehouseID == "" {
			sendError(log, w, http.StatusBadRequest, "warehouse_id is required")
			return
		}

		var req models.RequestCreateRoute
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		route, err := svc.CreateRoute(r.Context(), warehouseID, req)
		if err != nil {
			switch err {
			case coreErrors.ErrNotFoundWarehouse:
				sendError(log, w, http.StatusNotFound, "warehouse not found")
			case coreErrors.ErrDuplicateRoute:
				sendError(log, w, http.StatusConflict, "route already exists for this warehouse")
			default:
				sendError(log, w, http.StatusInternalServerError, "failed to create route")
			}
			return
		}

		sendCreated(log, w, route)
	}
}

func NewGetRoutesHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.PathValue("warehouse_id")
		if warehouseID == "" {
			sendError(log, w, http.StatusBadRequest, "warehouse_id is required")
			return
		}

		list, err := svc.GetRoutesByWarehouse(r.Context(), warehouseID)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to get routes")
			return
		}
		sendOK(log, w, list)
	}
}

func NewDeleteRouteHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.PathValue("warehouse_id")
		routeID := r.PathValue("route_id")
		if warehouseID == "" || routeID == "" {
			sendError(log, w, http.StatusBadRequest, "warehouse_id and route_id are required")
			return
		}

		if err := svc.DeleteRoute(r.Context(), warehouseID, routeID); err != nil {
			if err == coreErrors.ErrNotFoundRoute {
				sendError(log, w, http.StatusNotFound, "route not found")
				return
			}
			sendError(log, w, http.StatusInternalServerError, "failed to delete route")
			return
		}

		sendNoContent(w)
	}
}