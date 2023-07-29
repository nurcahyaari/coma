package entity

import "encoding/json"

type UserAccess struct {
	Id            string `json:"_id"`
	UserId        string `json:"userId"`
	ApplicationId string `json:"applicationId"`
	Stageid       string `json:"stageId"`
}

func (a UserAccess) MapStringInterface() (map[string]interface{}, error) {
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

type UserAccesses []UserAccess
