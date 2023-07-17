package localpubsub

import (
	"context"
	"io"
	"strings"

	"github.com/rs/zerolog/log"
)

func (h LocalPubsub) ConfigDistributor(id string, r io.Reader) {
	log.Info().
		Str("id", id).
		Msg("[ConfigDistributor] send configuration toward client")

	var clientKey string
	buf := new(strings.Builder)
	io.Copy(buf, r)
	clientKey = buf.String()

	err := h.configurationSvc.DistributeConfiguration(context.TODO(), clientKey)
	if err != nil {
		log.Error().Err(err).Msg("[ConfigDistributor] error distribute configuration")
		return
	}

	log.Info().
		Str("id", id).
		Msg("[ConfigDistributor] success send configuration toward client")
}
