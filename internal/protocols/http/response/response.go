package response

import (
	"encoding/json"
	"net/http"

	httperror "github.com/coma/coma/internal/protocols/http/errors"
)

type Response[T any] struct {
	Message *string `json:"message,omitempty"`
	Data    *T      `json:"data,omitempty"`
}

func Json[T any](w http.ResponseWriter, httpCode int, message string, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	res := Response[T]{
		Message: &message,
		Data:    &data,
	}
	json.NewEncoder(w).Encode(res)
}

func Text(w http.ResponseWriter, httpCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpCode)
	w.Write([]byte(message))
}

// TODO: implement response error
func Err[T any](w http.ResponseWriter, err error) {
	_, ok := err.(*httperror.RespError)
	if !ok {
		err = httperror.InternalServerError(err.Error())
	}

	er, _ := err.(*httperror.RespError)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(er.Code)
	res := Response[T]{
		Message: &er.Message,
	}
	json.NewEncoder(w).Encode(res)
}
