package model

import (
	"encoding/json"

	"github.com/ostafen/clover"
)

type ApplicationStage struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (r ApplicationStage) MapStringInterface() (map[string]interface{}, error) {
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

type ApplicationStages []ApplicationStage

type FilterApplicationStage struct {
	Name string
}

func (f FilterApplicationStage) Filter() *clover.Criteria {
	criterias := make([]*clover.Criteria, 0)

	if f.Name != "" {
		criterias = append(criterias, clover.Field("name").Eq(f.Name))
	}

	filter := &clover.Criteria{}

	for idx, criteria := range criterias {
		if idx == 0 {
			filter = criteria
			continue
		}

		filter = filter.And(criteria)
	}

	return filter
}
