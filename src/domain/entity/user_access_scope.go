package entity

import (
	"encoding/json"

	"github.com/ostafen/clover"
)

type UserApplicationScope struct {
	Id            string                    `json:"_id"`
	UserId        string                    `json:"userId"`
	ApplicationId string                    `json:"applicationId"`
	StageId       string                    `json:"stageId"`
	Rbac          *UserApplicationScopeRbac `json:"rbac"`
}

func (a *UserApplicationScope) UpdateRbac(userApplicationScopeNew UserApplicationScope) {
	if a.Rbac == nil && userApplicationScopeNew.Rbac == nil {
		return
	}

	a.Rbac = userApplicationScopeNew.Rbac
}

func (a UserApplicationScope) MapStringInterface() (map[string]interface{}, error) {
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

type UserApplicationsScope []UserApplicationScope

type UserApplicationScopeRbac struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type FilterUserApplicationScope struct {
	Id            string
	UserId        string
	UserIds       []string
	ApplicationId string
	StageId       string
}

func (f *FilterUserApplicationScope) Filter() *clover.Criteria {
	criterias := make([]*clover.Criteria, 0)

	if f.Id != "" {
		criterias = append(criterias, clover.Field("_id").Eq(f.Id))
	}

	if f.UserId != "" {
		criterias = append(criterias, clover.Field("userId").Eq(f.UserId))
	}

	if len(f.UserIds) > 0 {
		criterias = append(criterias, clover.Field("userId").In(f.UserIds))
	}

	if f.ApplicationId != "" {
		criterias = append(criterias, clover.Field("applicationId").Eq(f.ApplicationId))
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
