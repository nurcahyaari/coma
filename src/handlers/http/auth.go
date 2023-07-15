package http

import "net/http"

// SetConfig set new config
// @Summary set new config
// @Description Set new config
// @Tags Auth
// @Param Authorization header string true "<User Token>"
// @Produce json
// @Router /v1/auth/oauth/login [POST]
func (h *HttpHandle) OauthLogin(w http.ResponseWriter, r *http.Request) {

}
