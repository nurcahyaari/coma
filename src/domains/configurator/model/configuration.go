package model

type Configuration struct {
	Id    uint64 `json:"id"`
	Field string `json:"field"`
	Value string `json:"value"`
}
