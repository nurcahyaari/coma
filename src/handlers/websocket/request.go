package websocket

import (
	"encoding/json"
	"errors"
)

type RequestDistribute struct {
	ApiToken string          `json:"apiToken"`
	Data     json.RawMessage `json:"data"`
}

func (r RequestDistribute) Validate() []error {
	var errs []error

	if r.ApiToken == "" {
		errs = append(errs, errors.New("api token cannot be nulled or empty"))
	}

	if !json.Valid(r.Data) {
		errs = append(errs, errors.New("data must be a valid json"))
	}

	return errs
}
