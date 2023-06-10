package http

import "net/http"

// SetConfig set new config
// @Summary set new config
// @Description Set new config
// @Tags Config
// @Param Authorization header string true "<User Token>"
// @Produce json
// @Router /v1/configurator [POST]
func (h *HttpHandle) SetConfiguration(w http.ResponseWriter, r *http.Request) {

}
