package http

import (
	"net/http"

	"github.com/coma/coma/internal/protocols/http/response"
	"github.com/coma/coma/src/domains/application/dto"
)

func (h *HttpHandle) MiddlewareCheckIsClientKeyExists(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exist, _ := h.applicationKeySvc.IsExistsApplicationKey(r.Context(), dto.RequestFindApplicationKey{
			Key: r.Header.Get("x-clientkey"),
		})
		if !exist {
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
