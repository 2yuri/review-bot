package main

import (
	"github.com/2yuri/review-bot/api"
	"github.com/2yuri/review-bot/db"
	"github.com/2yuri/review-bot/db/repositories"
	"github.com/2yuri/review-bot/internal/sqs"
)

func main() {
	sqs.ConnectSQS()
	defer sqs.CloseConnectionSQS()

	client := db.Start()
	defer client.Close()

	uR := repositories.NewUserRepository(client)
	sR := repositories.NewSessionRepository(client)

	server := api.New("8080", uR, sR)
	server.Start()
}
