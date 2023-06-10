package dto

import (
	"encoding/json"
	"errors"
)

type RequestSetConfiguration struct {
	Data json.RawMessage `json:"data"`
}

func (r RequestSetConfiguration) Validate() error {
	if !json.Valid(r.Data) {
		return errors.New("err: not a valid json")
	}

	return nil
}
