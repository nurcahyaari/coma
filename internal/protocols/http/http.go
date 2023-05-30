package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coma/coma/config"
	"github.com/coma/coma/internal/protocols/http/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type HttpImpl struct {
	HttpRouter *router.HttpRouterImpl
	httpServer *http.Server
}

func NewHttp(httpRouter *router.HttpRouterImpl) *HttpImpl {
	return &HttpImpl{HttpRouter: httpRouter}
}

func (p *HttpImpl) setupRouter(app *chi.Mux) {
	p.HttpRouter.Router(app)
}

func (h *HttpImpl) cors(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
}

func (p *HttpImpl) Listen() {
	app := chi.NewRouter()

	p.setupRouter(app)

	serverPort := fmt.Sprintf(":%d", config.Get().Application.Port)
	p.httpServer = &http.Server{
		Addr:    serverPort,
		Handler: app,
	}

	log.Info().Msgf("Server started on Port %s ", serverPort)
	p.httpServer.ListenAndServe()
}

func (h *HttpImpl) Shutdown(ctx context.Context) error {
	if err := h.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
