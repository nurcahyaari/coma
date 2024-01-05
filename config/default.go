package config

import "time"

const (
	PUBSUB_DISTRIBUTOR_TOPIC = "coma:config:distributor:topic"
)

const (
	PUBSUB_DEFAULT_MAX_BUFFER_CAPACITY = 10000
	PUBSUB_DEFAULT_MAX_ELAPSED_TIME    = 10 * time.Millisecond
	PUBSUB_DEFAULT_MAX_WORKER          = 10000
	PUBSUB_DEFAULT_MAX_WAIT_TIME       = 25 * time.Millisecond
)
