package rest

import (
	"encoding/json"
	"net/http"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/ports"
)

func decodeAndValidate(log ports.Logger, w http.ResponseWriter, r *http.Request, v ports.Validator, dst any) bool {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		sendError(log, w, http.StatusBadRequest, "invalid json: "+err.Error())
		return false
	}
	if err := v.ValidateStruct(dst); err != nil {
		sendError(log, w, http.StatusBadRequest, coreErrors.ErrValidateFailed.Error())
		return false
	}
	return true
}