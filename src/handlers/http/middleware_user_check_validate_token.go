package http

import (
	"net/http"
	"strings"

	"github.com/coma/coma/internal/protocols/http/response"
	"github.com/coma/coma/src/application/auth/dto"
	"github.com/coma/coma/src/domain/entity"
)

func (h *HttpHandle) MiddlewareLocalAuthAccessTokenValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")

		token := strings.Split(authorization, " ")

		if len(token) < 2 && token[0] != string(entity.BearerAuthenticationToken) {
			response.Err[string](
				w,
				response.SetErr[string]("err: token mismatch"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}

		resp, err := h.authSvc.ValidateToken(r.Context(), dto.RequestValidateToken{
			Token:     token[1],
			TokenType: entity.AccessToken,
		})
		if err != nil {
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}
		if !resp.Valid {
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
