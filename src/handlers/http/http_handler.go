package http

import (
	applicationsvc "github.com/coma/coma/src/domains/application/service"
	authsvc "github.com/coma/coma/src/domains/auth/service"
	"github.com/go-chi/chi/v5"
)

type HttpHandle struct {
	authSvc             authsvc.Servicer
	configurationSvc    applicationsvc.ApplicationConfigurationServicer
	applicationStageSvc applicationsvc.ApplicationStageServicer
	applicationSvc      applicationsvc.ApplicationServicer
	applicationKeySvc   applicationsvc.ApplicationKeyServicer
}

func (h HttpHandle) Router(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Route("/oauth", func(r chi.Router) {
				r.Post("/login", h.OauthLogin)
			})
		})

		r.Route("/applications", func(r chi.Router) {
			r.Get("/", h.FindApplications)
			r.Post("/", h.CreateApplication)
			r.Delete("/{applicationId}", h.DeleteApplications)
		})

		r.Route("/stages", func(r chi.Router) {
			r.Get("/", h.FindApplicationStages)
			r.Post("/", h.CreateApplicationStages)
			r.Delete("/{stageName}", h.DeleteApplicationStages)
		})

		r.Route("/keys", func(r chi.Router) {
			r.Get("/", h.FindApplicationKey)
			r.Post("/", h.CreateOrUpdateApplicationKey)
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
	configuratorSvc applicationsvc.ApplicationConfigurationServicer,
	applicationEnvSvc applicationsvc.ApplicationStageServicer,
	applicationSvc applicationsvc.ApplicationServicer,
	applicationKeySvc applicationsvc.ApplicationKeyServicer) HttpOption {
	return func(h *HttpHandle) {
		h.authSvc = authSvc
		h.configurationSvc = configuratorSvc
		h.applicationStageSvc = applicationEnvSvc
		h.applicationSvc = applicationSvc
		h.applicationKeySvc = applicationKeySvc
	}
}

func NewHttpHandler(opts ...HttpOption) *HttpHandle {
	httpHandle := &HttpHandle{}

	for _, opt := range opts {
		opt(httpHandle)
	}

	return httpHandle
}
