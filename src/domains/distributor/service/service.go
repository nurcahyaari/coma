package service

import (
	"encoding/json"

	"github.com/coma/coma/src/domains/distributor/dto"
	selfextdto "github.com/coma/coma/src/external/self/dto"
	selfextsvc "github.com/coma/coma/src/external/self/service"
)

type Servicer interface {
	SendMessage() error
}

type Service struct {
	selfExtSvc selfextsvc.WSServicer
}

func New(distributor selfextsvc.WSServicer) Servicer {
	return &Service{
		selfExtSvc: distributor,
	}
}

func (s *Service) SendMessage() error {
	messageJson := dto.RequestDistribute{
		ApiToken: "12345",
		Data:     `"{\n  \"apiToken\": \"123456\",\n  \"data\": {\n    \"port\": \"1234\"\n  },\n  \"contentType\": \"json\"\n}"`,
	}

	vJson, _ := json.Marshal(messageJson)

	s.selfExtSvc.Send(selfextdto.RequestSendMessage{
		Message: vJson,
	})

	return nil
}
