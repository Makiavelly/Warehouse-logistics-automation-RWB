package rest

import (
	"net/http"
	"time"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

func NewGetTruckCallsHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.URL.Query().Get("warehouse_id")
		routeID := r.URL.Query().Get("route_id")

		var from, to time.Time
		if s := r.URL.Query().Get("from"); s != "" {
			from, _ = time.Parse(time.RFC3339, s)
		}
		if s := r.URL.Query().Get("to"); s != "" {
			to, _ = time.Parse(time.RFC3339, s)
		}

		list, err := svc.GetTruckCalls(r.Context(), warehouseID, routeID, from, to)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to get truck calls")
			return
		}
		sendOK(log, w, list)
	}
}

func NewGetTruckCallAccuracyHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		warehouseID := r.URL.Query().Get("warehouse_id")
		routeID := r.URL.Query().Get("route_id")

		acc, err := svc.GetTruckCallAccuracy(r.Context(), warehouseID, routeID)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to get accuracy")
			return
		}
		sendOK(log, w, acc)
	}
}

func NewReportTimelinessHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		truckCallID := r.PathValue("truck_call_id")
		if truckCallID == "" {
			sendError(log, w, http.StatusBadRequest, "truck_call_id is required")
			return
		}

		var req models.RequestReportTimeliness
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		tc, err := svc.ReportTimeliness(r.Context(), truckCallID, req)
		if err != nil {
			if err == coreErrors.ErrNotFoundTruckCall {
				sendError(log, w, http.StatusNotFound, "truck call not found")
				return
			}
			sendError(log, w, http.StatusInternalServerError, "failed to report timeliness")
			return
		}

		sendOK(log, w, tc)
	}
}