package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/coma/coma/internal/protocols/http/response"
	"github.com/coma/coma/src/domains/auth/service"
	"github.com/go-chi/chi/v5"
)

type HttpHandlerImpl struct {
	svc service.Servicer
}

func (h HttpHandlerImpl) Router(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello test")
		time.Sleep(10 * time.Second)
		response.Text(w, 200, "hello")
	})
}
func NewHttpHandler(svc service.Servicer) *HttpHandlerImpl {
	return &HttpHandlerImpl{
		svc: svc,
	}
}
