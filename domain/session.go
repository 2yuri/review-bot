package domain

import (
	"fmt"
	"github.com/2yuri/review-bot/internal/sqs"
	"github.com/2yuri/review-bot/internal/telegram"
)

type SessionStatus string

const (
	Waiting    SessionStatus = "WAITING"
	Processing SessionStatus = "PROCESSING"
	Done       SessionStatus = "DONE"
)

type SessionType string

const (
	Review SessionType = "REVIEW"
	Info   SessionType = "INFO"
)

type Session struct {
	ID         int64         `json:"-"`
	ExternalId string        `json:"externalId"`
	Name       string        `json:"name"`
	Product    string        `json:"product"`
	ChatID     string        `json:"-"`
	Status     SessionStatus `json:"status"`
	Type       SessionType   `json:"type"`
	Logs       []SessionLogs `json:"logs"`
}

type SessionLogs struct {
	ID   int64  `json:"-"`
	Text string `json:"text"`
}

func (s *Session) FormatMessage() sqs.Message {
	switch s.Type {
	case Review:
		return sqs.Message{
			Data:    fmt.Sprintf(telegram.ReviewMessage, s.Product),
			ChatID:  s.ChatID,
			Buttons: []string{"1 ☹️", "2 😕", "3 😐", "4 🙃", "5 🤩"},
		}
	}

	return sqs.Message{}
}
