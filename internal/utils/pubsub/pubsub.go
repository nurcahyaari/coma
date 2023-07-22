package pubsub

import (
	"context"
	"errors"
	"log"

	"github.com/coma/coma/internal/utils/pubsub/database"
	"github.com/ostafen/clover"
)

var (
	ErrTopicIsNotExists    error = errors.New("err: topic is not exists")
	ErrConsumerIsNotExists error = errors.New("err: your message buffer already full but you don't have any consumer yet")
)

type Pubsub struct {
	shutdown   chan bool
	database   database.Databaser
	publisher  map[string]*publisher
	subscriber map[string]*subscriber
}

type PubsubOption func(pb *Pubsub)

func SetCloverForBackup(db *clover.DB) PubsubOption {
	return func(pb *Pubsub) {
		database := database.Database{
			DatabaseDriver: database.CLOVER,
		}
		pb.database = database.NewCloverDatabase(db)
	}
}

func NewPubsub(opts ...PubsubOption) *Pubsub {
	pubsub := &Pubsub{
		shutdown:   make(chan bool),
		publisher:  make(map[string]*publisher),
		subscriber: make(map[string]*subscriber),
	}

	for _, opt := range opts {
		opt(pubsub)
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

	// check is there backup for this topic
	go ps.CheckBackup(topic)
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

		if subscriber.handler == nil {
			continue
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

	log.Printf("publish message to %s, current message: %d\n", topic, pub.len())

	return pub.publish(message)
}

func (ps Pubsub) Capacity(topic string) int {
	return ps.publisher[topic].capacity()
}

func (ps Pubsub) Len(topic string) int {
	return ps.publisher[topic].len()
}

func (ps Pubsub) CheckBackup(topic string) error {
	if ps.database == nil {
		return nil
	}
	log.Println("checking backup...")
	backups, err := ps.database.RetrieveAndDelete(topic)
	if err != nil {
		return err
	}
	if len(backups) == 0 {
		log.Println("no backup...")
		return nil
	}

	log.Printf("found %d backups, start publish...\n", len(backups))

	for _, backup := range backups {
		if err := ps.Publish(backup.Topic, SendBytes(backup.Message)); err != nil {
			return err
		}
	}

	log.Printf("success retrieve %d message...\n", len(backups))
	return nil
}

func (ps Pubsub) shutdownSubscriber(topic string) {
	ps.subscriber[topic].close()
}

func (ps Pubsub) shutdownPublisher(topic string) {
	ps.publisher[topic].close()
}

func (ps Pubsub) Shutdown(ctx context.Context) error {
	// todo: implement later
	if ps.database == nil {
		return nil
	}
	log.Println("backup message from queue")

	backups := database.Backups{}
	for topic, publisher := range ps.publisher {
		ps.shutdownSubscriber(topic)
		ps.shutdownPublisher(topic)
		messages, err := publisher.retrieveMessages()
		if err != nil {
			return err
		}

		for _, message := range messages {
			backups = append(backups, database.Backup{
				Topic:   topic,
				Message: []byte(message),
			})
		}
	}

	if len(backups) == 0 {
		log.Println("there is no message from queue")
		return nil
	}

	for _, backup := range backups {
		err := ps.database.Store(backup)
		if err != nil {
			return err
		}
	}

	log.Printf("success backup %d message from queue\n", len(backups))
	return nil
}
