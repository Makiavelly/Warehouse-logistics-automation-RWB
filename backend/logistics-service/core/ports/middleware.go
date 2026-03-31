package ports

import "net/http"

type Middleware interface {
	Auth(tokenGen TokenGenerator, role Role) func(http.HandlerFunc) http.HandlerFunc
	Logging(next http.Handler) http.Handler
}