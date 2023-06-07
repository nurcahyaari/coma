package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/coma/coma/config"
	"github.com/coma/coma/internal/graceful"
	"github.com/coma/coma/internal/protocols/http/response"
	"github.com/coma/coma/internal/protocols/http/router"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"
)

type HttpImpl struct {
	HttpRouter  *router.HttpRouterImpl
	httpServer  *http.Server
	serverState graceful.ServerState
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

func (h *HttpImpl) shutdownStateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server state", h.serverState)
		switch h.serverState {
		case graceful.StateShutdown:
			response.Json[string](w, http.StatusInternalServerError, "server is shutting down", "")
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}

func (p *HttpImpl) Listen() {
	app := chi.NewRouter()

	p.cors(app)
	app.Use(p.shutdownStateMiddleware)
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
	h.serverState = graceful.StateShutdown
	if err := h.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	h.HttpRouter.CloseWebsocket()

	return nil
}
