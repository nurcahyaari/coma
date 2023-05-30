package http

import (
	"github.com/coma/coma/src/domains/auth/service"
	"github.com/go-chi/chi/v5"
)

type HttpHandlerImpl struct {
	svc service.Servicer
}

func (h HttpHandlerImpl) Router(r *chi.Mux) {
}
func NewHttpHandler(svc service.Servicer) *HttpHandlerImpl {
	return &HttpHandlerImpl{
		svc: svc,
	}
}
