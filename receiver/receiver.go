package receiver

import (
	"bytes"
	"compress/gzip"
	"io"
	"time"

	"github.com/pebbe/zmq4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func ReceiveData(exit chan bool) {
	subscriber, _ := zmq4.NewSocket(zmq4.SUB)

	defer subscriber.Close()

	subscriber.SetLinger(0)
	subscriber.SetRcvtimeo(1 * time.Second)

	zmqHost := viper.GetString("ndov.server")

	initializeSources()

	subscriber.Connect(zmqHost)
	log.WithField("host", zmqHost).Info("Connect to server")

	subscribeSources(subscriber)

	listen(subscriber, exit)
}

func listen(subscriber *zmq4.Socket, exit chan bool) {
	log.Info("Receiving data...")

	for {
		select {
		case <-exit:
			log.Info("Shutting down receiver")

			subscriber.Close()
			log.Info("Receiver shut down")

			exit <- true

			return
		default:
			msg, err := subscriber.RecvMessageBytes(0)

			if err != nil {
				continue
			}

			envelope := string(msg[0])

			// Decompress message:
			message, _ := gunzip(msg[1])

			if err != nil {
				log.WithFields(log.Fields{
					"error":    err,
					"envelope": envelope,
					"message":  string(msg[1]),
				}).Error("Error decompressing message. Message ignored")
			} else {
				source, found := lookupSource(envelope)

				if !found {
					log.WithFields(log.Fields{
						"envelope": envelope,
					}).Warning("Unknown envelope")
					log.Debug(message)
				} else {
					log.WithFields(log.Fields{
						"source": source,
					}).Info("Process message")
				}
			}
		}
	}
}

func gunzip(data []byte) (io.Reader, error) {
	buf := bytes.NewBuffer(data)
	reader, err := gzip.NewReader(buf)

	if err != nil {
		reader.Close()
		return nil, err
	}

	defer reader.Close()

	buf3 := new(bytes.Buffer)
	buf3.ReadFrom(reader)

	return buf3, nil
}
