package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	l, _ := zap.NewProduction()
	logger = l
	defer logger.Sync()
}

type Event struct {
	Name string `json:"name"`
}

func MyHandler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		logger.Info("received sqs event", zap.Any("message", message))

		event := &Event{}
		err := json.Unmarshal([]byte(message.Body), event)
		if err != nil {
			return err
		}

		logger.Info("received sqs event", zap.Any("event", event))
	}

	return nil
}

func main() {
	lambda.Start(MyHandler)
}
