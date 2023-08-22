package http

import (
	"fmt"
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	"github.com/coma/coma/src/application/auth/dto"
	"github.com/rs/zerolog/log"
)

func (h *HttpHandle) MiddlewareLocalAuthUserScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := h.authSvc.ValidateUserScope(r.Context(), dto.RequestUserScopeValidation{
			UserId: r.Header.Get("x-coma-user-id"),
			Method: r.Method,
		})
		if err != nil {
			log.Error().
				Msg("[MiddlewareLocalAuthUserScope.ValidateUserScope] err")
			response.Err[string](
				w,
				response.SetErr[string]("err: forbidden"),
				response.SetHttpCode[string](http.StatusForbidden),
			)
			return
		}

		if !resp.Valid {
			log.Error().
				Str("user", r.Header.Get("x-coma-user-id")).
				Str("method", r.Method).
				Msg("[MiddlewareLocalAuthUserScope.ValidateUserScope] forbidden")
			response.Err[string](
				w,
				response.SetErr[string]("err: forbidden"),
				response.SetHttpCode[string](http.StatusForbidden),
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// TODO: add middleware based on user access scope
func (h *HttpHandle) MiddlewareLocalAuthUserAccessScope(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method)
		next.ServeHTTP(w, r)
	})
}
