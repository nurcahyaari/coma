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

func ReaderToReaderSeeker(reader io.Reader) (io.ReadSeeker, error) {
	// Create a buffer to hold the entire content
	var buffer bytes.Buffer

	// Copy the content from the reader to the buffer
	_, err := io.Copy(&buffer, reader)
	if err != nil {
		return nil, err
	}

	// Create a ReadSeeker from the buffer
	readSeeker := bytes.NewReader(buffer.Bytes())

	return readSeeker, nil
}
