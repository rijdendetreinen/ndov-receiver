package output

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

type RedisOutput struct {
	config OutputConfig
	rdb    *redis.Client
}

func (s *RedisOutput) Config() OutputConfig {
	return s.config
}

func (redisOutput *RedisOutput) Setup(config OutputConfig) {
	redisOutput.config = config
	opt, err := redis.ParseURL(config.URL)

	if err != nil {
		panic(err)
	}

	redisOutput.rdb = redis.NewClient(opt)
}

func (output *RedisOutput) ProcessMessage(source string, message string) {
	log.WithField("source", source).WithField("output", output.config.Name).Debug("Processing message")

	redisResult := output.rdb.LPush(ctx, source, message)

	if redisResult.Err() != nil {
		log.WithField("source", source).WithField("output", output.config.Name).WithError(redisResult.Err()).Error("Error while processing data")
	}
}
