package receiver

import (
	"strings"

	"github.com/pebbe/zmq4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var sourcesEnvelopes map[string]string

func initializeSources() {
	sourcesEnvelopes = viper.GetStringMapString("ndov.sources")
}

func subscribeSources(subscriber *zmq4.Socket) {
	for key, envelope := range sourcesEnvelopes {
		log.WithFields(log.Fields{
			"source":   key,
			"envelope": envelope,
		}).Info("Subscribed to source")
		subscriber.SetSubscribe(envelope)
	}
}

func lookupSource(envelope string) (string, bool) {
	for source, sourcesEnvelope := range sourcesEnvelopes {
		if strings.HasPrefix(envelope, sourcesEnvelope) {
			return source, true
		}
	}

	return "", false
}
