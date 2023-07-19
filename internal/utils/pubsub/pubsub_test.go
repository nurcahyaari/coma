package pubsub_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/coma/coma/internal/utils/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestPubsub(t *testing.T) {
	t.Run("test consumer string", func(t *testing.T) {
		actual := make(chan string)
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			buf := new(strings.Builder)
			io.Copy(buf, r)
			actual <- buf.String()
			close(actual)
		}, pubsub.PubsubSetMaxWorker(5))

		go ps.Listen()

		ps.Publish("test-topic-1", pubsub.SendString("hello world"))

		assert.Equal(t, "hello world", <-actual)
	})

	t.Run("test consumer bytes", func(t *testing.T) {
		actual := make(chan string)
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			buf := new(strings.Builder)
			io.Copy(buf, r)
			actual <- buf.String()
			close(actual)
		}, pubsub.PubsubSetMaxWorker(5))

		go ps.Listen()

		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))

		assert.Equal(t, "hello bytes", <-actual)
	})

	t.Run("test consumer json", func(t *testing.T) {
		type JSON struct {
			Message string `json:"message"`
		}
		actual := make(chan JSON)
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			var resp JSON
			json.NewDecoder(r).Decode(&resp)
			actual <- resp
			close(actual)
		}, pubsub.PubsubSetMaxWorker(5))

		go ps.Listen()

		text := "{\"message\":\"Hello Json\"}"
		ps.Publish("test-topic-1", pubsub.SendJSON(json.RawMessage(text)))

		expected := JSON{
			Message: "Hello Json",
		}
		assert.Equal(t, expected, <-actual)
	})

	t.Run("test publisher capacity", func(t *testing.T) {
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		assert.Equal(t, 5, ps.Capacity("test-topic-1"))
	})

	t.Run("test publisher len", func(t *testing.T) {
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		assert.Equal(t, 0, ps.Len("test-topic-1"))
	})
}

func TestShutdown(t *testing.T) {
	t.Run("message queue is empty", func(t *testing.T) {
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			fmt.Println(id, r)
		}, pubsub.PubsubSetMaxWorker(5))

		go ps.Listen()

		err := ps.Shutdown(context.TODO())
		assert.NoError(t, err)
	})
}
