package pubsub

import "io"

type publisher struct {
	message chan io.Reader
}

type publisherOptions struct {
	bufferCapacity int
}

func newPublisher(options publisherOptions) *publisher {
	pub := &publisher{
		message: make(chan io.Reader, options.bufferCapacity),
	}

	return pub
}

func (p *publisher) publish(message MessageHandler) error {
	data, err := message()
	if err != nil {
		return err
	}

	p.message <- data
	return nil
}

func (p *publisher) capacity() int {
	return cap(p.message)
}

func (p *publisher) len() int {
	return len(p.message)
}
