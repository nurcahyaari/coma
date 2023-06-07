package service

import (
	"encoding/json"

	"github.com/coma/coma/src/domains/distributor/dto"
	distributorextdto "github.com/coma/coma/src/external/distributor/dto"
	distributorextsvc "github.com/coma/coma/src/external/distributor/service"
)

type Servicer interface {
	SendMessage() error
}

type Service struct {
	distributor distributorextsvc.WSServicer
}

func New(distributor distributorextsvc.WSServicer) Servicer {
	return &Service{
		distributor: distributor,
	}
}

func (s *Service) SendMessage() error {
	messageJson := dto.RequestDistribute{
		ApiToken: "12345",
		Data:     `"{\n  \"apiToken\": \"123456\",\n  \"data\": {\n    \"port\": \"1234\"\n  },\n  \"contentType\": \"json\"\n}"`,
	}

	vJson, _ := json.Marshal(messageJson)

	s.distributor.Send(distributorextdto.RequestSendMessage{
		Message: vJson,
	})

	return nil
}
