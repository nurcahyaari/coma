package database

import "encoding/json"

type Backup struct {
	Id         string `json:"_id"`
	SequenceId int64  `json:"sequenceId"`
	Topic      string `json:"topic"`
	Message    []byte `json:"byte"`
}

func (r Backup) MapStringInterface() (map[string]interface{}, error) {
	mapStringIntf := make(map[string]interface{})
	j, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &mapStringIntf)
	if err != nil {
		return nil, err
	}
	return mapStringIntf, nil
}

type Backups []Backup
