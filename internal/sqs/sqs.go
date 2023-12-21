package sqs

import (
	"context"
	"github.com/2yuri/review-bot/internal/env"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"log"
)

var sqsCl *sqs.Client

type Message struct {
	Data    string   `json:"data"`
	ChatID  string   `json:"chatId"`
	Buttons []string `json:"buttons"`
}

func ConnectSQS() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(env.Get("REGION").StringFallback("us-east-1")))
	if err != nil {
		log.Fatalln("cannot load config: ", err)
	}

	sqsCl = sqs.NewFromConfig(cfg)
}

func CloseConnectionSQS() {
	sqsCl = nil
}

func SendToQueue(ctx context.Context, queue, message string) error {
	_, err := sqsCl.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &queue,
	})

	return err
}
