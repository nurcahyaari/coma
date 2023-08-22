package http

import (
	"net/http"
	"strings"

	"github.com/coma/coma/internal/protocols/http/response"
	"github.com/coma/coma/src/application/auth/dto"
	"github.com/coma/coma/src/domain/entity"
	"github.com/rs/zerolog/log"
)

func (h *HttpHandle) MiddlewareLocalAuthAccessTokenValidate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("authorization")

		token := strings.Split(authorization, " ")

		if len(token) < 2 && token[0] != string(entity.BearerAuthenticationToken) {
			log.Error().
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ValidateToken] token mismatch")
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
			log.Error().
				Err(err).
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ValidateToken]")
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}
		if !resp.Valid {
			log.Error().
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ValidateToken] token isn't valid")
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}

		extracted, err := h.authSvc.ExtractToken(r.Context(), dto.RequestValidateToken{
			Token:     token[1],
			TokenType: entity.AccessToken,
		})
		if err != nil {
			log.Error().
				Msg("[MiddlewareLocalAuthAccessTokenValidate.ExtractToken] token cannot be extracted")
			response.Err[string](
				w,
				response.SetErr[string]("err: unauthorized"),
				response.SetHttpCode[string](http.StatusUnauthorized),
			)
			return
		}

		r.Header.Set("x-coma-user-id", extracted.UserId)
		r.Header.Set("x-coma-user-type", extracted.UserType)

		next.ServeHTTP(w, r)
	})
}
