package localpubsub

import (
	"github.com/coma/coma/config"
	"github.com/coma/coma/container"
	"github.com/coma/coma/internal/utils/pubsub"
	"github.com/coma/coma/src/domains/service"
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
	h.pubSub.ConsumerRegister(h.config.Pubsub.Local.Consumer.ConfigDistributor.Topic, h.ConfigDistributor)
}

func (h LocalPubsub) TopicRegistry() {
	h.pubSub.TopicRegister(h.config.Pubsub.Local.Publisher.ConfigDistributor.Topic,
		pubsub.PubsubSetMaxBufferCapacity(h.config.Pubsub.Local.Publisher.ConfigDistributor.MaxBufferCapacity))
}

func (h LocalPubsub) Listen() {
	h.Consumer()

	go h.pubSub.Listen()
}
