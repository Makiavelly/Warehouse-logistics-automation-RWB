package middleware

import (
	"net/http"
	"time"
)

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.log.Debug("request started", "method", r.Method, "path", r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r)
		m.log.Debug("request done", "method", r.Method, "path", r.URL.Path, "elapsed", time.Since(start))
	})
}