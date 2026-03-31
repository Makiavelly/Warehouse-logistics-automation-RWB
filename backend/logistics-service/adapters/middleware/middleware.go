package middleware

import "logistics-service/logistics-service/core/ports"

type Middleware struct {
	log    ports.Logger
	apiKey string
}

func NewMiddleware(log ports.Logger, apiKey string) *Middleware {
	return &Middleware{log: log, apiKey: apiKey}
}