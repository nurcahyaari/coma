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
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog/log"

	"net/http/pprof"

	_ "github.com/coma/coma/docs"
	httpswagger "github.com/swaggo/http-swagger"
)

type Http struct {
	HttpRouter  *router.HttpRoute
	httpServer  *http.Server
	serverState graceful.ServerState
}

func NewHttp(httpRouter *router.HttpRoute) *Http {
	return &Http{HttpRouter: httpRouter}
}

func (p *Http) setupRouter(app *chi.Mux) {
	p.HttpRouter.Router(app)
}

func (h *Http) cors(r *chi.Mux) {
	r.Use(cors.AllowAll().Handler)
}

func (h *Http) setupSwagger(app *chi.Mux) {
	app.Mount("/swagger", httpswagger.WrapHandler)
	app.Handle("/public/*", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
}

func (h *Http) shutdownStateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch h.serverState {
		case graceful.StateShutdown:
			// response.Json[string](w, http.StatusInternalServerError, "server is shutting down", "")
			response.Json[string](w,
				response.SetHttpCode[string](http.StatusInternalServerError),
				response.SetMessage[string]("server is shutting down"))
			return
		default:
			next.ServeHTTP(w, r)
		}
	})
}

func (h *Http) setupPprof(r *chi.Mux) {
	r.HandleFunc("/debug/pprof/", pprof.Index)
	r.HandleFunc("/debug/pprof/heap", pprof.Index)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func (p *Http) Listen() {
	app := chi.NewRouter()

	p.cors(app)
	app.Use(middleware.Logger)
	app.Use(middleware.Recoverer)
	app.Use(p.shutdownStateMiddleware)
	p.setupRouter(app)
	p.setupSwagger(app)
	p.setupPprof(app)

	serverPort := fmt.Sprintf(":%d", config.Get().Application.Port)
	p.httpServer = &http.Server{
		Addr:    serverPort,
		Handler: app,
	}

	log.Info().Msgf("Server started on Port %s ", serverPort)
	p.httpServer.ListenAndServe()
}

func (h *Http) Shutdown(ctx context.Context) error {
	h.serverState = graceful.StateShutdown
	if err := h.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	h.HttpRouter.CloseWebsocket()

	return nil
}
