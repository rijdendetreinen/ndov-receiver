package output

import (
	"context"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type RedisOutput struct {
	output Output
	rdb    *redis.Client
}

var redisOutputs map[string]RedisOutput = make(map[string]RedisOutput)

var ctx = context.Background()

func setupRedisOutput(output Output) {
	opt, err := redis.ParseURL(output.URL)

	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)

	redisOutput := RedisOutput{
		output: output,
		rdb:    rdb,
	}

	redisOutputs[output.Name] = redisOutput
}

func getRedisOutput(output Output) RedisOutput {
	return redisOutputs[output.Name]
}

func processRedisMessage(output Output, source string, message string) {
	redisOutput := getRedisOutput(output)

	log.WithField("source", source).WithField("output", output.Name).Debug("Processing message")

	redisResult := redisOutput.rdb.LPush(ctx, source, message)

	if redisResult.Err() != nil {
		log.WithField("source", source).WithField("output", output.Name).WithError(redisResult.Err()).Error("Error while processing data")
	}
}
