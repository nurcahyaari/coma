package dto

import "github.com/coma/coma/src/domains/application/model"

type RequestDeleteConfiguration struct {
	XClientKey string
	Id         string
}

func (r RequestDeleteConfiguration) FilterConfiguration() model.FilterConfiguration {
	return model.FilterConfiguration{
		Id:        r.Id,
		ClientKey: r.XClientKey,
	}
}
