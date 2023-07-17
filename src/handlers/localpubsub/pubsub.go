package localpubsub

import (
	"github.com/coma/coma/config"
	"github.com/coma/coma/internal/utils/pubsub"
	applicationsvc "github.com/coma/coma/src/domains/application/service"
)

type LocalPubsub struct {
	config           *config.Config
	pubSub           *pubsub.Pubsub
	configurationSvc applicationsvc.ApplicationConfigurationServicer
}

type LocalPubsubOption func(h *LocalPubsub)

func SetDomains(
	configurationSvc applicationsvc.ApplicationConfigurationServicer) LocalPubsubOption {
	return func(h *LocalPubsub) {
		h.configurationSvc = configurationSvc
	}
}

func NewLocalPubsub(config *config.Config, pubSub *pubsub.Pubsub, opts ...LocalPubsubOption) *LocalPubsub {
	localPubsub := &LocalPubsub{
		config: config,
		pubSub: pubSub,
	}

	for _, opt := range opts {
		opt(localPubsub)
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
