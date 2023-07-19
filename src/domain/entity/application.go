package entity

import (
	"encoding/json"

	"github.com/ostafen/clover"
)

type Application struct {
	Id      string `json:"_id"`
	StageId string `json:"stageId"`
	Type    string `json:"type"`
	Name    string `json:"name"`
}

func (a Application) MapStringInterface() (map[string]interface{}, error) {
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

type Applications []Application

type FilterApplication struct {
	Id      string
	Name    string
	StageId string
}

func (f FilterApplication) Filter() *clover.Criteria {
	criterias := make([]*clover.Criteria, 0)

	if f.Id != "" {
		criterias = append(criterias, clover.Field("_id").Eq(f.Id))
	}

	if f.Name != "" {
		criterias = append(criterias, clover.Field("name").Eq(f.Name))
	}

	if f.StageId != "" {
		criterias = append(criterias, clover.Field("stageId").Eq(f.StageId))
	}

	filter := &clover.Criteria{}

	if len(criterias) == 0 {
		return nil
	}

	for idx, criteria := range criterias {
		if idx == 0 {
			filter = criteria
			continue
		}

		filter = filter.And(criteria)
	}

	return filter
}
