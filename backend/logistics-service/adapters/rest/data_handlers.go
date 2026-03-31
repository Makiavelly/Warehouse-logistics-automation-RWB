package rest

import (
	"encoding/json"
	"net/http"

	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

func NewIngestDataHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestIngestData
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		resp, err := svc.IngestData(r.Context(), req)
		if err != nil {
			sendError(log, w, http.StatusInternalServerError, "failed to ingest data")
			return
		}
		sendOK(log, w, resp)
	}
}

func NewRetrainModelHandler(log ports.Logger, svc ports.Service, _ ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// body is optional
		var req models.RequestRetrainModel
		_ = json.NewDecoder(r.Body).Decode(&req)

		resp, err := svc.RetrainModel(r.Context(), req)
		if err != nil {
			sendError(log, w, http.StatusServiceUnavailable, "retrain failed: "+err.Error())
			return
		}
		sendOK(log, w, resp)
	}
}