package http

import (
	applicationsvc "github.com/coma/coma/src/domains/application/service"
	authsvc "github.com/coma/coma/src/domains/auth/service"
	configuratorsvc "github.com/coma/coma/src/domains/configurator/service"
	"github.com/go-chi/chi/v5"
)

type HttpHandle struct {
	authSvc             authsvc.Servicer
	configurationSvc    configuratorsvc.Servicer
	applicationStageSvc applicationsvc.ApplicationStageServicer
}

func (h HttpHandle) Router(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Route("/oauth", func(r chi.Router) {
				r.Post("/login", h.OauthLogin)
			})
		})

		r.Route("/stages", func(r chi.Router) {
			r.Post("/", h.CreateApplicationStages)
		})

		r.Route("/configurator", func(r chi.Router) {
			r.Get("/", h.GetConfiguration)
			r.Post("/", h.SetConfiguration)
			r.Put("/", h.UpdateConfiguration)
			r.Post("/upsert", h.UpsertConfiguration)
			r.Delete("/{id}", h.DeleteConfiguration)
		})
	})
}

type HttpOption func(h *HttpHandle)

func SetDomains(
	authSvc authsvc.Servicer,
	configuratorSvc configuratorsvc.Servicer,
	applicationEnvSvc applicationsvc.ApplicationStageServicer) HttpOption {
	return func(h *HttpHandle) {
		h.authSvc = authSvc
		h.configurationSvc = configuratorSvc
		h.applicationStageSvc = applicationEnvSvc
	}
}

func NewHttpHandler(opts ...HttpOption) *HttpHandle {
	httpHandle := &HttpHandle{}

	for _, opt := range opts {
		opt(httpHandle)
	}

	return httpHandle
}
