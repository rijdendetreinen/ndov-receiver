package output

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type OutputConfig struct {
	Name       string
	OutputType string
	URL        string
	Sources    []string
}

type Output interface {
	Setup(config OutputConfig)
	Config() OutputConfig
	ProcessMessage(source string, message string)
}

var outputs map[string]Output = make(map[string]Output)

func SetupOutputs() {
	log.Debug("Setting up outputs")

	outputConfigs := viper.GetStringMap("outputs")

	for outputName := range outputConfigs {
		config := OutputConfig{
			Name:       outputName,
			OutputType: viper.GetString("outputs." + outputName + ".type"),
			URL:        viper.GetString("outputs." + outputName + ".url"),
			Sources:    viper.GetStringSlice("outputs." + outputName + ".sources"),
		}

		switch config.OutputType {
		case "redis":
			log.WithField("output", outputName).WithField("type", config.OutputType).Info("Setting up output")

			redisOutput := &RedisOutput{}
			redisOutput.Setup(config)

			outputs[outputName] = redisOutput

		case "file":
			log.WithField("output", outputName).WithField("type", config.OutputType).Info("Setting up output")

			fileOutput := &FileOutput{}
			fileOutput.Setup(config)

			outputs[outputName] = fileOutput

		default:
			log.WithField("output", outputName).WithField("type", config.OutputType).Error("Invalid output type")
		}
	}
}

func ProcessMessage(source string, message string) {
	for _, output := range outputs {
		if !contains(output.Config().Sources, source) {
			continue
		}

		output.ProcessMessage(source, message)
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
