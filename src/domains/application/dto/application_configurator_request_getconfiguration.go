package dto

type RequestGetConfiguration struct {
	XClientKey string `json:"clientKey"`
}

const (
	ViewTypeSchema = "schema"
	ViewTypeJSON   = "JSON"
)
