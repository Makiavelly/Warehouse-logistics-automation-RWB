package rest

import (
	"net/http"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
	"logistics-service/logistics-service/adapters/middleware"
)

func NewGetDriversHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := svc.GetDrivers(r.Context())
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to get drivers")
			return
		}
		sendOK(log, w, list)
	}
}

func NewAssignDriverHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestAssignDriver
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		driver, err := svc.AssignDriver(r.Context(), req)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to assign driver")
			return
		}
		sendOK(log, w, driver)
	}
}

func NewGetDriverSignalHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.UserIDFromCtx(r.Context())

		tc, err := svc.GetDriverSignal(r.Context(), userID)
		if err != nil {
			if err == coreErrors.ErrNotFoundDriver {
				sendError(log, w, http.StatusNotFound, "driver profile not found")
				return
			}
			sendError(log, w, http.StatusInternalServerError, "failed to get signal")
			return
		}

		if tc == nil {
			sendOK(log, w, map[string]any{"signal": false, "truck_call": nil})
			return
		}
		sendOK(log, w, map[string]any{"signal": true, "truck_call": tc})
	}
}

func NewGetDriverStatsHandler(log ports.Logger, svc ports.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.UserIDFromCtx(r.Context())

		stats, err := svc.GetDriverStats(r.Context(), userID)
		if err != nil {
			if err == coreErrors.ErrNotFoundDriver {
				sendError(log, w, http.StatusNotFound, "driver profile not found")
				return
			}
			sendError(log, w, http.StatusInternalServerError, "failed to get stats")
			return
		}
		sendOK(log, w, stats)
	}
}