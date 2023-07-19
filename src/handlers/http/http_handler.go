package http

import (
	service "github.com/coma/coma/src/domains/service"
	"github.com/go-chi/chi/v5"
)

type HttpHandle struct {
	authSvc             service.AuthServicer
	configurationSvc    service.ApplicationConfigurationServicer
	applicationStageSvc service.ApplicationStageServicer
	applicationSvc      service.ApplicationServicer
	applicationKeySvc   service.ApplicationKeyServicer
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

		r.Route("/configuration", func(r chi.Router) {
			r.Use(h.MiddlewareCheckIsClientKeyExists)
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
	authSvc service.AuthServicer,
	configurationSvc service.ApplicationConfigurationServicer,
	applicationEnvSvc service.ApplicationStageServicer,
	service service.ApplicationServicer,
	applicationKeySvc service.ApplicationKeyServicer) HttpOption {
	return func(h *HttpHandle) {
		h.authSvc = authSvc
		h.configurationSvc = configurationSvc
		h.applicationStageSvc = applicationEnvSvc
		h.applicationSvc = service
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
