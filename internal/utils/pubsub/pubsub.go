package pubsub

import (
	"context"
	"errors"
)

var (
	ErrTopicIsNotExists    error = errors.New("err: topic is not exists")
	ErrConsumerIsNotExists error = errors.New("err: your message buffer already full but you don't have any consumer yet")
)

type Pubsub struct {
	publisher  map[string]*publisher
	subscriber map[string]*subscriber
}

func NewPubsub() *Pubsub {
	pubsub := &Pubsub{
		publisher:  make(map[string]*publisher),
		subscriber: make(map[string]*subscriber),
	}

	return pubsub
}

type PubsubRegisterOpt func(ps *PubsubRegisterOptions)

type PubsubRegisterOptions struct {
	maxBufferCapacity int
}

func PubsubSetMaxBufferCapacity(max int) PubsubRegisterOpt {
	return func(ps *PubsubRegisterOptions) {
		ps.maxBufferCapacity = max
	}
}

// register the callback
func (ps Pubsub) TopicRegister(topic string, opts ...PubsubRegisterOpt) {
	var pubsubRegisterOption PubsubRegisterOptions

	for _, opt := range opts {
		opt(&pubsubRegisterOption)
	}

	ps.publisher[topic] = newPublisher(publisherOptions{
		bufferCapacity: pubsubRegisterOption.maxBufferCapacity,
	})

	ps.subscriber[topic] = newSubscriber()
}

func (ps Pubsub) ConsumerRegister(topic string, handler SubscriberHandler, opts ...SubscriberOption) error {
	_, exists := ps.subscriber[topic]
	if !exists {
		return ErrTopicIsNotExists
	}
	ps.subscriber[topic].registerSubscriberHandler(handler, opts...)
	return nil
}

// consume the message
func (ps Pubsub) Listen() error {
	for topic, publisher := range ps.publisher {
		subscriber, exists := ps.subscriber[topic]
		if !exists || subscriber == nil {
			return ErrTopicIsNotExists
		}

		go subscriber.dispatcher(publisher)
		go subscriber.listen()
	}

	return nil
}

func (ps Pubsub) Publish(topic string, message MessageHandler) error {
	pub, exists := ps.publisher[topic]
	if !exists {
		return ErrTopicIsNotExists
	}

	if _, exists := ps.subscriber[topic]; !exists {
		return ErrConsumerIsNotExists
	}

	return pub.publish(message)
}

func (ps Pubsub) Capacity(topic string) int {
	return ps.publisher[topic].capacity()
}

func (ps Pubsub) Len(topic string) int {
	return ps.publisher[topic].len()
}

func (ps Pubsub) Shutdown(ctx context.Context) error {
	// todo: implement later
	return nil
}
