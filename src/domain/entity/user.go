package entity

import "encoding/json"

type User struct {
	Id       string `json:"_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a User) MapStringInterface() (map[string]interface{}, error) {
	mapStringIntf := make(map[string]interface{})
	j, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &mapStringIntf)
	if err != nil {
		return nil, err
	}
	return mapStringIntf, nil
}

type Users []User

type FilterUser struct {
}
