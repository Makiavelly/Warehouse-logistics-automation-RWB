package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	coreErrors "logistics-service/logistics-service/core/errors"
	"logistics-service/logistics-service/core/ports"
)

type contextKey string

const (
	CtxUserID = contextKey("user_id")
	CtxRole   = contextKey("role")
)

// Auth middleware checks JWT and requires the given role ("admin" or "driver").
// Pass role="" to allow any authenticated user.
func (m *Middleware) Auth(
	tokenGen ports.TokenGenerator,
	roleRequired ports.Role,
) func(http.HandlerFunc) http.HandlerFunc {
	const op = "middleware.Auth"

	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenStr := extractToken(r)
			if tokenStr == "" {
				writeAuthError(w, m.log, "token is empty")
				return
			}

			claims, err := tokenGen.Validate(tokenStr)
			if err != nil {
				m.log.Warn(op+" token validation failed", "error", err)
				writeAuthError(w, m.log, "invalid token")
				return
			}

			if roleRequired != "" && claims.Role != roleRequired {
				writeAuthError(w, m.log, "insufficient permissions")
				return
			}

			ctx := context.WithValue(r.Context(), CtxUserID, claims.UserID)
			ctx = context.WithValue(ctx, CtxRole, claims.Role)
			h.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

// APIKeyAuth middleware checks a static API key header for external integrations.
func (m *Middleware) APIKeyAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")
		if key == "" || key != m.apiKey {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(coreErrors.NewErrResponse("invalid api key"))
			return
		}
		h.ServeHTTP(w, r)
	}
}

func extractToken(r *http.Request) string {
	auth := r.Header.Get("Authorization")
	if strings.HasPrefix(auth, "Bearer ") {
		return strings.TrimPrefix(auth, "Bearer ")
	}
	return r.Header.Get("token")
}

func writeAuthError(w http.ResponseWriter, log ports.Logger, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(coreErrors.NewErrResponse(msg)); err != nil {
		log.Debug("writeAuthError encode failed", "error", err)
	}
}

func UserIDFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(CtxUserID).(string)
	return v
}

func RoleFromCtx(ctx context.Context) string {
	v, _ := ctx.Value(CtxRole).(string)
	return v
}
