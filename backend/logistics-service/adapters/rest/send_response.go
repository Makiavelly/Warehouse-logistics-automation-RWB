package rest

import (
	"encoding/json"
	"net/http"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/ports"
)

func sendOK(log ports.Logger, w http.ResponseWriter, body any) {
	sendJSON(log, w, http.StatusOK, body)
}

func sendCreated(log ports.Logger, w http.ResponseWriter, body any) {
	sendJSON(log, w, http.StatusCreated, body)
}

func sendNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func sendJSON(log ports.Logger, w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Debug("sendJSON encode error", "error", err)
	}
}

func sendError(log ports.Logger, w http.ResponseWriter, status int, msg string) {
	sendJSON(log, w, status, coreErrors.NewErrResponse(msg))
}