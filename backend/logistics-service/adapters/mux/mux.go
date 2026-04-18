package mux

import "net/http"

type Mux struct {
	*http.ServeMux
}

func NewMux() *Mux {
	return &Mux{http.NewServeMux()}
}