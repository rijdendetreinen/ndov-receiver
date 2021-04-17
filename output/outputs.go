package output

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Output struct {
	Name       string
	OutputType string
	URL        string
	Sources    []string
}

var outputs map[string]Output = make(map[string]Output)

func SetupOutputs() {
	log.Debug("Setting up outputs")

	outputConfigs := viper.GetStringMap("outputs")

	for outputName, _ := range outputConfigs {
		output := Output{
			Name:       outputName,
			OutputType: viper.GetString("outputs." + outputName + ".type"),
			URL:        viper.GetString("outputs." + outputName + ".url"),
			Sources:    viper.GetStringSlice("outputs." + outputName + ".sources"),
		}

		switch output.OutputType {
		case "redis":
			log.WithField("output", outputName).WithField("type", output.OutputType).Info("Setting up output")

			setupRedisOutput(output)

			outputs[outputName] = output

		default:
			log.WithField("output", outputName).WithField("type", output.OutputType).Error("Invalid output type")
		}
	}
}

func ProcessMessage(source string, message string) {
	for _, output := range outputs {
		if !contains(output.Sources, source) {
			continue
		}

		switch output.OutputType {
		case "redis":
			processRedisMessage(output, source, message)
		}
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
