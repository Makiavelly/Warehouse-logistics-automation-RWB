package rest

import (
	"net/http"
	"time"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

func NewRequestForecastHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestPredict
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		forecast, err := svc.RequestForecast(r.Context(), req)
		if err != nil {
			switch err {
			case coreErrors.ErrNotFoundWarehouse:
				sendError(log, w, http.StatusNotFound, "warehouse not found")
			case coreErrors.ErrNotFoundRoute:
				sendError(log, w, http.StatusNotFound, "route not found")
			case coreErrors.ErrMLServiceUnavailable:
				sendError(log, w, http.StatusServiceUnavailable, "ml service unavailable")
			default:
				sendError(log, w, http.StatusInternalServerError, "failed to get forecast")
			}
			return
		}

		sendCreated(log, w, forecast)
	}
}

func NewGetForecastsHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := models.RequestForecastQuery{
			WarehouseID: r.URL.Query().Get("warehouse_id"),
			RouteID:     r.URL.Query().Get("route_id"),
		}

		if fromStr := r.URL.Query().Get("from"); fromStr != "" {
			if t, err := time.Parse(time.RFC3339, fromStr); err == nil {
				q.From = t
			}
		}
		if toStr := r.URL.Query().Get("to"); toStr != "" {
			if t, err := time.Parse(time.RFC3339, toStr); err == nil {
				q.To = t
			}
		}

		list, err := svc.GetForecasts(r.Context(), q)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to get forecasts")
			return
		}
		sendOK(log, w, list)
	}
}