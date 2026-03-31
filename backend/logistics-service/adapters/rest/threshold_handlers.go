package rest

import (
	"net/http"

	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

func NewSetThresholdHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestSetThreshold
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		threshold, err := svc.SetThreshold(r.Context(), req)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to set threshold")
			return
		}

		sendOK(log, w, threshold)
	}
}

func NewGetThresholdsHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.URL.Query().Get("warehouse_id")
		routeID := r.URL.Query().Get("route_id")

		list, err := svc.GetThresholds(r.Context(), warehouseID, routeID)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to get thresholds")
			return
		}
		sendOK(log, w, list)
	}
}