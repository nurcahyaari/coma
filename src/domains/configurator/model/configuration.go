package model

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/ostafen/clover"
	"gopkg.in/guregu/null.v4"
)

type Configuration struct {
	Id        string `json:"id"`
	ClientKey string `json:"clientKey"`
	Field     string `json:"field"`
	Value     any    `json:"value"`
}

func (c Configuration) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Id, validation.Required),
		validation.Field(&c.ClientKey, validation.Required),
		validation.Field(&c.Field, validation.Required),
		validation.Field(&c.Value, validation.Required),
	)
}

func (r *Configuration) Update(configuration Configuration) {
	r.ClientKey = configuration.ClientKey
	r.Field = configuration.Field
	r.Value = configuration.Value
}

func (r Configuration) MapStringInterface() (map[string]interface{}, error) {
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

func (r Configuration) FilterConfiguration() FilterConfiguration {
	return FilterConfiguration{
		Id:        r.Id,
		ClientKey: r.ClientKey,
	}
}

type Configurations []Configuration

func (rs Configurations) Exists() bool {
	return len(rs) > 0
}

func (rs Configurations) Update(mapConfiguration MapConfigurationById) {
	for idx, r := range rs {
		newConfiguration, ok := mapConfiguration[r.Id]
		if !ok {
			continue
		}
		rs[idx].Update(newConfiguration)
	}
}

func (rs Configurations) MapFieldValue() (map[string]any, error) {
	mapFieldValue := make(map[string]any)

	for _, r := range rs {
		mapFieldValue[r.Field] = r.Value
	}

	return mapFieldValue, nil
}

func (rs Configurations) MapConfigurationById() MapConfigurationById {
	mapConfigurationById := make(MapConfigurationById)
	for _, r := range rs {
		mapConfigurationById[r.Id] = r
	}
	return mapConfigurationById
}

type MapConfigurationById map[string]Configuration

// FilterConfiguration lets you filter its data, the argument is "and"
type FilterConfiguration struct {
	Id          string
	ParentField null.String
	ClientKey   string
	Field       string
}

func (f FilterConfiguration) Filter() *clover.Criteria {
	criterias := make([]*clover.Criteria, 0)

	if f.ClientKey != "" {
		criterias = append(criterias, clover.Field("clientKey").Eq(f.ClientKey))
	}

	if f.Id != "" {
		criterias = append(criterias, clover.Field("id").Eq(f.Id))
	}

	if f.Field != "" {
		criterias = append(criterias, clover.Field("field").Eq(f.Field))
	}

	if f.ParentField.Valid {
		criterias = append(criterias, clover.Field("parentField").Eq(f.ParentField))
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
