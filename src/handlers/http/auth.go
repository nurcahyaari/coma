package http

import (
	"net/http"

	"github.com/coma/coma/src/application/auth/dto"
)

// SetConfig set new config
// @Summary set new config
// @Description Set new config
// @Tags Auth
// @Param Authorization header string true "<User Token>"
// @Produce json
// @Router /v1/auth/oauth/login [POST]
func (h *HttpHandle) OauthLogin(w http.ResponseWriter, r *http.Request) {
	h.authSvc.ValidateToken(r.Context(), dto.RequestValidateToken{
		Method: dto.Oauth,
		Token:  "hello",
	})
}
