package http

import (
	"github.com/coma/coma/container"
	service "github.com/coma/coma/src/domain/service"
	"github.com/go-chi/chi/v5"
)

type HttpHandle struct {
	authSvc             service.LocalUserAuthServicer
	configurationSvc    service.ApplicationConfigurationServicer
	applicationStageSvc service.ApplicationStageServicer
	applicationSvc      service.ApplicationServicer
	applicationKeySvc   service.ApplicationKeyServicer
	userSvc             service.UserServicer
}

func (h HttpHandle) Router(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/applications", func(r chi.Router) {
			r.Use(
				h.MiddlewareLocalAuthAccessTokenValidate,
				h.MiddlewareLocalAuthUserApplicationScope,
				h.MiddlewareLocalAuthUserScope)
			r.Get("/", h.FindApplications)
			r.Post("/", h.CreateApplication)
			r.Delete("/{applicationId}", h.DeleteApplications)
		})

		r.Route("/stages", func(r chi.Router) {
			r.Use(
				h.MiddlewareLocalAuthAccessTokenValidate,
				h.MiddlewareLocalAuthUserApplicationScope,
				h.MiddlewareLocalAuthUserScope)
			r.Get("/", h.FindApplicationStages)
			r.Post("/", h.CreateApplicationStages)
			r.Delete("/{stageName}", h.DeleteApplicationStages)
		})

		r.Route("/keys", func(r chi.Router) {
			r.Use(
				h.MiddlewareLocalAuthAccessTokenValidate,
				h.MiddlewareLocalAuthUserApplicationScope,
				h.MiddlewareLocalAuthUserScope)
			r.Get("/", h.FindApplicationKey)
			r.Post("/", h.CreateOrUpdateApplicationKey)
		})

		r.Route("/configuration", func(r chi.Router) {
			r.Use(
				h.MiddlewareLocalAuthAccessTokenValidate,
				h.MiddlewareCheckIsClientKeyExists,
				h.MiddlewareLocalAuthUserApplicationScope,
				h.MiddlewareLocalAuthUserScope)
			r.Get("/", h.GetConfiguration)
			r.Post("/", h.SetConfiguration)
			r.Put("/", h.UpdateConfiguration)
			r.Post("/upsert", h.UpsertConfiguration)
			r.Delete("/{id}", h.DeleteConfiguration)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/root", h.CreateUserRoot)
			r.Group(func(r chi.Router) {
				r.Use(h.MiddlewareLocalAuthAccessTokenValidate)
				r.Patch("/password/{id}", h.UpdateUserPassword)
				r.Group(func(r chi.Router) {
					r.Use(h.MiddlewareLocalAuthUserScope)
					r.Get("/", h.FindUsers)
					r.Get("/{id}", h.FindUser)
					r.Post("/", h.CreateUser)
					r.Delete("/{id}", h.DeleteUser)
					r.Put("/{id}", h.UpdateUser)
				})
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Post("/login", h.AuthUserLogin)
			})
		})
	})
}

func NewHttpHandler(c container.Service) *HttpHandle {
	httpHandle := &HttpHandle{
		authSvc:             c.LocalUserAuthServicer,
		configurationSvc:    c.ApplicationConfigurationServicer,
		applicationStageSvc: c.ApplicationStageServicer,
		applicationSvc:      c.ApplicationServicer,
		applicationKeySvc:   c.ApplicationKeyServicer,
		userSvc:             c.UserServicer,
	}
	return httpHandle
}
