package pubsub

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type MessageHandler func() (io.Reader, error)

func SendString(data string) MessageHandler {
	return func() (io.Reader, error) {
		return strings.NewReader(data), nil
	}
}

func SendBytes(data []byte) MessageHandler {
	return func() (io.Reader, error) {
		return bytes.NewBuffer(data), nil
	}
}

func SendJSON(data json.RawMessage) MessageHandler {
	return func() (io.Reader, error) {
		var buff bytes.Buffer

		err := json.NewEncoder(&buff).Encode(data)
		if err != nil {
			return nil, err
		}

		return &buff, nil
	}
}
