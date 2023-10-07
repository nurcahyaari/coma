package pubsub

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"

	"github.com/coma/coma/internal/utils/pubsub/database"
	"github.com/ostafen/clover"
)

var (
	ErrTopicIsNotExists    error = errors.New("err: topic is not exists")
	ErrConsumerIsNotExists error = errors.New("err: your message buffer already full but you don't have any consumer yet")
)

type Pubsub struct {
	shutdown          chan bool
	database          database.Databaser
	publisher         map[string]*publisher
	subscriber        map[string][]*subscriber
	subscriberCounter int
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
		subscriber: make(map[string][]*subscriber),
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
func (ps *Pubsub) TopicRegister(topic string, opts ...PubsubRegisterOpt) {
	var pubsubRegisterOption PubsubRegisterOptions

	for _, opt := range opts {
		opt(&pubsubRegisterOption)
	}

	ps.publisher[topic] = newPublisher(publisherOptions{
		bufferCapacity: pubsubRegisterOption.maxBufferCapacity,
	})

	// check is there backup for this topic
	go ps.CheckBackup(topic)
}

func (ps *Pubsub) ConsumerRegister(topic string, handler SubscriberHandler, opts ...SubscriberOption) error {
	defer func() {
		ps.subscriberCounter++
	}()

	newSubscriber := newSubscriber()
	newSubscriber.registerSubscriberHandler(handler, opts...)
	ps.subscriber[topic] = append(ps.subscriber[topic], newSubscriber)

	go ps.dispatcher(topic)
	return nil
}

// consume the message
func (ps *Pubsub) Listen() error {
	for topic, _ := range ps.publisher {
		subscribers, exists := ps.subscriber[topic]
		if !exists || subscribers == nil || len(subscribers) == 0 {
			return ErrTopicIsNotExists
		}

		go ps.listener(topic)
	}

	return nil
}

func (ps Pubsub) dispatcher(topic string) {
	for {
		select {
		case <-ps.shutdown:
			return
		case message := <-ps.publisher[topic].message:
			buff := new(bytes.Buffer)
			_, err := io.Copy(buff, message)
			if err != nil {
				continue
			}

			for _, subscriber := range ps.subscriber[topic] {
				if subscriber.handler == nil {
					continue
				}

				duplicateMessage := bytes.NewBufferString(buff.String())

				go subscriber.dispatcher(duplicateMessage)
			}
		}
	}
}

func (ps Pubsub) listener(topic string) {
	for _, subscriber := range ps.subscriber[topic] {
		if subscriber.handler == nil {
			continue
		}

		go subscriber.listen()
	}
}

func (ps *Pubsub) Publish(topic string, message MessageHandler) error {
	pub, exists := ps.publisher[topic]
	if !exists {
		return ErrTopicIsNotExists
	}

	if err := pub.publish(message); err != nil {
		return err
	}

	if _, exists := ps.subscriber[topic]; !exists {
		log.Printf("topic %s doesn have subscriber the message will store to the memory, current message: %d\n", topic, pub.len())
		return ErrConsumerIsNotExists
	}

	log.Printf("success publishing the message to %s, current message: %d\n", topic, pub.len())

	return nil
}

func (ps Pubsub) Capacity(topic string) int {
	return ps.publisher[topic].capacity()
}

func (ps Pubsub) Len(topic string) int {
	return ps.publisher[topic].len()
}

func (ps *Pubsub) CheckBackup(topic string) error {
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

func (ps *Pubsub) shutdownSubscriber(topic string) {
	for idx, _ := range ps.subscriber[topic] {
		ps.subscriber[topic][idx].close()
	}
}

func (ps *Pubsub) shutdownPublisher(topic string) {
	ps.publisher[topic].close()
}

func (ps *Pubsub) Shutdown(ctx context.Context) error {
	if ps.database == nil {
		log.Println("no database was selected")
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
