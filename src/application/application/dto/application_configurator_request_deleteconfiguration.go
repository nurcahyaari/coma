package dto

import "github.com/coma/coma/src/domains/entity"

type RequestDeleteConfiguration struct {
	XClientKey string
	Id         string
}

func (r RequestDeleteConfiguration) FilterConfiguration() entity.FilterConfiguration {
	return entity.FilterConfiguration{
		Id:        r.Id,
		ClientKey: r.XClientKey,
	}
}
