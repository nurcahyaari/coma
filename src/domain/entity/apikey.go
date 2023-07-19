package entity

type Apikey struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key"`
}
