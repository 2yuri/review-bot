package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/2yuri/review-bot/db"
	"github.com/2yuri/review-bot/db/repositories"
	"github.com/2yuri/review-bot/domain"
	"github.com/2yuri/review-bot/internal/env"
	l "github.com/2yuri/review-bot/internal/logger"
	"github.com/2yuri/review-bot/internal/sqs"
	"go.uber.org/zap"
	"log"
)

var Repository repositories.SessionRepository

func main() {
	sqs.ConnectSQS()
	defer sqs.CloseConnectionSQS()
	//lambda.Start(handler)

	ctx := context.Background()

	client := db.Start()
	defer client.Close()

	Repository = repositories.NewSessionRepository(client)

	fmt.Println("running manually: ", handler(ctx))
}

func handler(ctx context.Context) error {
	reviews, err := Repository.GetPendingSessions(ctx)
	if err != nil {
		return err
	}

	if len(reviews) == 0 {
		log.Println("nothing to compute")
		return nil
	}

	for _, review := range reviews {
		err := Repository.SetStatus(ctx, domain.Processing, review.ID)
		if err != nil {
			log.Println("cannot set review status: ", err)
			return err
		}

		b, err := json.Marshal(review.FormatMessage())
		if err != nil {
			l.Logger.Error("cannot marshal message", zap.Error(err))

			return err
		}

		if err := sqs.SendToQueue(ctx, env.Get("DISPATCHER_QUEUE_URL").String(), string(b)); err != nil {
			return err
		}
	}

	return nil
}
