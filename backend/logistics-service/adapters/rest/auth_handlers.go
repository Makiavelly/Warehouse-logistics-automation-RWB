package rest

import (
	"net/http"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/models"
	"logistics-service/logistics-service/core/ports"
)

func NewLoginHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestLogin
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		resp, err := svc.Login(r.Context(), req)
		if err != nil {
			sendError(log, w, http.StatusUnauthorized, coreErrors.ErrInvalidCredentials.Error())
			return
		}

		sendOK(log, w, resp)
	}
}

func NewRegisterHandler(log ports.Logger, svc ports.Service, v ports.Validator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.RequestRegister
		if !decodeAndValidate(log, w, r, v, &req) {
			return
		}

		user, err := svc.Register(r.Context(), req)
		if err != nil {
			if err == coreErrors.ErrDuplicateUser {
				sendError(log, w, http.StatusConflict, "username already exists")
				return
			}
			sendError(log, w, http.StatusInternalServerError, "failed to register user")
			return
		}

		sendCreated(log, w, user)
	}
}
