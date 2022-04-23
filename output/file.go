package output

import (
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

type FileOutput struct {
	config OutputConfig
}

func (s *FileOutput) Config() OutputConfig {
	return s.config
}

func (output *FileOutput) Setup(config OutputConfig) {
	output.config = config
}

func (output *FileOutput) ProcessMessage(source string, message string) {
	now := time.Now()

	filename := output.config.URL + "/" + source + "." + fmt.Sprint(now.UnixNano()) + ".log"

	log.WithField("filename", filename).WithField("source", source).WithField("output", output.config.Name).Debug("Processing message")
	err := ioutil.WriteFile(filename, []byte(message), 0644)

	if err != nil {
		log.Error(err)
	}
}
