package localpubsub

import (
	"github.com/nurcahyaari/coma/config"
	"github.com/nurcahyaari/coma/container"
	"github.com/nurcahyaari/coma/internal/x/pubsub"
	"github.com/nurcahyaari/coma/src/domain/service"
)

type LocalPubsub struct {
	config           *config.Config
	pubSub           *pubsub.Pubsub
	configurationSvc service.ApplicationConfigurationServicer
}

func NewLocalPubsub(config *config.Config, c container.Container) *LocalPubsub {
	localPubsub := &LocalPubsub{
		config:           config,
		pubSub:           c.LocalPubsub,
		configurationSvc: c.ApplicationConfigurationServicer,
	}
	return localPubsub
}

func (h LocalPubsub) Consumer() {
	h.pubSub.ConsumerRegister(h.config.Pubsub.ConfigDistributor.Consumer.Topic, h.ConfigDistributor)
}

func (h LocalPubsub) TopicRegistry() {
	h.pubSub.TopicRegister(h.config.Pubsub.ConfigDistributor.Publisher.Topic,
		pubsub.PubsubSetMaxBufferCapacity(h.config.Pubsub.ConfigDistributor.Publisher.MaxBufferCapacity))
}

func (h LocalPubsub) Listen() {
	h.Consumer()

	go h.pubSub.Listen()
}
