package websocket

import (
	"encoding/json"
	"errors"
)

type RequestDistribute struct {
	ClientKey string          `json:"clientKey"`
	Data      json.RawMessage `json:"data"`
}

func (r RequestDistribute) Validate() []error {
	var errs []error

	if r.ClientKey == "" {
		errs = append(errs, errors.New("client key cannot be nulled or empty"))
	}

	if !json.Valid(r.Data) {
		errs = append(errs, errors.New("data must be a valid json"))
	}

	return errs
}
