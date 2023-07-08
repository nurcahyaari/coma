package response

import (
	"encoding/json"
	"net/http"
)

type Response[T any, E any] struct {
	Message *string `json:"message,omitempty"`
	Data    *T      `json:"data,omitempty"`
	Err     *E      `json:"error,omitempty"`
}

type ResponseData[T any] struct {
	HttpCode int
	Message  string
	Data     T
	Err      T
}

type ResponseOption[T any] func(r *ResponseData[T])

func SetHttpCode[T any](httpCode int) ResponseOption[T] {
	return func(r *ResponseData[T]) {
		r.HttpCode = httpCode
	}
}

func SetMessage[T any](message string) ResponseOption[T] {
	return func(r *ResponseData[T]) {
		r.Message = message
	}
}

func SetData[T any](data T) ResponseOption[T] {
	return func(r *ResponseData[T]) {
		r.Data = data
	}
}

func SetErr[T any](err T) ResponseOption[T] {
	return func(r *ResponseData[T]) {
		r.Err = err
	}
}

func Json[T any](w http.ResponseWriter, opts ...ResponseOption[T]) {
	respData := &ResponseData[T]{
		HttpCode: 200, // default http code
	}

	for _, opt := range opts {
		opt(respData)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respData.HttpCode)
	res := Response[T, string]{
		Message: &respData.Message,
		Data:    &respData.Data,
	}

	json.NewEncoder(w).Encode(res)
}

func Text(w http.ResponseWriter, httpCode int, message string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(httpCode)
	w.Write([]byte(message))
}

func Err[E any](w http.ResponseWriter, opts ...ResponseOption[E]) {
	respData := &ResponseData[E]{}

	for _, opt := range opts {
		opt(respData)
	}

	// if user didn't set httpcode properly
	// it automatically settled as internal error
	if respData.HttpCode <= http.StatusOK {
		respData.HttpCode = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(respData.HttpCode)
	res := Response[any, E]{
		Message: &respData.Message,
		Err:     &respData.Err,
	}
	json.NewEncoder(w).Encode(res)
}
