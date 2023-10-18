package pubsub_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	"github.com/coma/coma/internal/x/pubsub"
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

	t.Run("test multi consumer string", func(t *testing.T) {
		actual := make(chan string, 2)
		wg := sync.WaitGroup{}
		wg.Add(2)
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			defer wg.Done()
			fmt.Println(r, id, "consumer 1")
			buf1 := new(strings.Builder)
			io.Copy(buf1, r)
			actual <- buf1.String()
		}, pubsub.PubsubSetMaxWorker(1))

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			defer wg.Done()
			fmt.Println(r, id, "consumer 2")
			buf2 := new(strings.Builder)
			io.Copy(buf2, r)
			actual <- buf2.String()
		}, pubsub.PubsubSetMaxWorker(2))

		go ps.Listen()

		ps.Publish("test-topic-1", pubsub.SendString("hello world"))

		wg.Wait()
		close(actual)

		exp := "hello world"
		counter := 1
		for val := range actual {
			assert.Equal(t, exp, val, fmt.Sprintf("consumer %d", counter))
			counter++
		}
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

	t.Run("test multi consumer json", func(t *testing.T) {
		type JSON struct {
			Message string `json:"message"`
		}

		actual := make(chan JSON, 2)
		wg := sync.WaitGroup{}
		wg.Add(2)

		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(5),
		)

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			defer wg.Done()
			var resp JSON
			json.NewDecoder(r).Decode(&resp)
			actual <- resp
		}, pubsub.PubsubSetMaxWorker(5))

		ps.ConsumerRegister("test-topic-1", func(id string, r io.Reader) {
			defer wg.Done()
			var resp JSON
			json.NewDecoder(r).Decode(&resp)
			actual <- resp
		}, pubsub.PubsubSetMaxWorker(2))

		go ps.Listen()

		text := "{\"message\":\"Hello Json\"}"
		ps.Publish("test-topic-1", pubsub.SendJSON(json.RawMessage(text)))

		expected := JSON{
			Message: "Hello Json",
		}

		wg.Wait()
		close(actual)

		counter := 1
		for val := range actual {
			assert.Equal(t, expected, val, fmt.Sprintf("consumer %d", counter))
			counter++
		}
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

	t.Run("test publisher len more than 0", func(t *testing.T) {
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(100),
		)

		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))

		go ps.Listen()

		assert.Equal(t, 5, ps.Len("test-topic-1"))
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

	t.Run("after shutdown it should backup the message", func(t *testing.T) {
		ps := pubsub.NewPubsub()
		ps.TopicRegister(
			"test-topic-1",
			pubsub.PubsubSetMaxBufferCapacity(100),
		)

		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))
		ps.Publish("test-topic-1", pubsub.SendBytes([]byte("hello bytes")))

		go ps.Listen()

		err := ps.Shutdown(context.TODO())
		assert.NoError(t, err)
	})

}
