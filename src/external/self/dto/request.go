package dto

import "encoding/json"

type RequestSendMessage struct {
	ClientKey string          `json:"clientKey"`
	Data      json.RawMessage `json:"data"`
}

func (r RequestSendMessage) Message() ([]byte, error) {
	byt, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return byt, nil
}
