package main

import (
	"context"
	"encoding/json"
	l "github.com/2yuri/review-bot/internal/logger"
	"github.com/2yuri/review-bot/internal/sqs"
	"github.com/2yuri/review-bot/internal/telegram"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.SQSEvent) error {
	var messages []sqs.Message

	for _, record := range event.Records {
		var message sqs.Message
		err := json.Unmarshal([]byte(record.Body), &message)
		if err != nil {
			continue
		}

		messages = append(messages, message)
	}

	l.Logger.Debug("Send Telegram", zap.Any("messages", messages))

	for _, m := range messages {
		err := telegram.SendBotMessage(m.ChatID, m.Data, m.Buttons)
		if err != nil {
			return err
		}
	}

	return nil
}
