package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/nurcahyaari/coma/container"
	service "github.com/nurcahyaari/coma/src/domain/service"
)

type HttpHandle struct {
	authSvc                 service.LocalUserAuthServicer
	configurationSvc        service.ApplicationConfigurationServicer
	applicationSvc          service.ApplicationServicer
	applicationKeySvc       service.ApplicationKeyServicer
	userSvc                 service.UserServicer
	userApplicationScopeSvc service.UserApplicationScopeServicer
}

func (h HttpHandle) Router(r *chi.Mux) {
	r.Route("/v1", func(r chi.Router) {
		r.Route("/applications", func(r chi.Router) {
			r.Use(
				h.MiddlewareLocalAuthAccessTokenValidate,
				// h.MiddlewareLocalAuthUserApplicationScope, TODO: uncomment later
				h.MiddlewareLocalAuthUserScope)
			r.Get("/", h.FindApplications)
			r.Post("/", h.CreateApplication)
			r.Delete("/{applicationId}", h.DeleteApplications)
		})

		r.Route("/keys", func(r chi.Router) {
			r.Use(
				h.MiddlewareLocalAuthAccessTokenValidate,
				// h.MiddlewareLocalAuthUserApplicationScope, TODO: uncomment later
				h.MiddlewareLocalAuthUserScope)
			r.Get("/", h.FindApplicationKey)
			r.Post("/", h.CreateOrUpdateApplicationKey)
		})

		r.Route("/configuration", func(r chi.Router) {
			r.Use(
				h.MiddlewareLocalAuthAccessTokenValidate,
				h.MiddlewareCheckIsClientKeyExists,
				// h.MiddlewareLocalAuthUserApplicationScope, TODO: uncomment later
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
					r.Post("/", h.CreateUser)
					r.Route("/{id}", func(r chi.Router) {
						r.Get("/", h.FindUser)
						r.Delete("/", h.DeleteUser)
						r.Put("/", h.UpdateUser)
					})
					r.Route("/application", func(r chi.Router) {
						r.Get("/scope", h.FindUserApplicationsScope)
						r.Post("/scope", h.CreateOrUpdateUserApplicationScope)
					})
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
		authSvc:                 c.LocalUserAuthServicer,
		configurationSvc:        c.ApplicationConfigurationServicer,
		applicationSvc:          c.ApplicationServicer,
		applicationKeySvc:       c.ApplicationKeyServicer,
		userSvc:                 c.UserServicer,
		userApplicationScopeSvc: c.UserApplicationScopeServicer,
	}
	return httpHandle
}
