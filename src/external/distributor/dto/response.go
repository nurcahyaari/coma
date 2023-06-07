package dto

import "encoding/json"

type ResponseStringData struct {
	Message string
}

type ResponseJSONData struct {
	Message json.RawMessage
}
