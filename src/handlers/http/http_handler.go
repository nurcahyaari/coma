package http

import (
	authsvc "github.com/coma/coma/src/domains/auth/service"
	configuratorsvc "github.com/coma/coma/src/domains/configurator/service"
	"github.com/go-chi/chi/v5"
)

type HttpHandle struct {
	authSvc          authsvc.Servicer
	configurationSvc configuratorsvc.Servicer
}

func (h HttpHandle) Router(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Route("/oauth", func(r chi.Router) {
				r.Post("/login", h.OauthLogin)
			})
		})
		r.Route("/configurator", func(r chi.Router) {
			r.Post("/", h.SetConfiguration)
		})
	})
}

type HttpOption func(h *HttpHandle)

func SetDomains(authSvc authsvc.Servicer, configuratorSvc configuratorsvc.Servicer) HttpOption {
	return func(h *HttpHandle) {
		h.authSvc = authSvc
		h.configurationSvc = configuratorSvc
	}
}

func NewHttpHandler(opts ...HttpOption) *HttpHandle {
	httpHandle := &HttpHandle{}

	for _, opt := range opts {
		opt(httpHandle)
	}

	return httpHandle
}
