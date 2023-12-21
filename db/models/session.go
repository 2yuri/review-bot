package models

import (
	"github.com/2yuri/review-bot/domain"
	"time"
)

type Session struct {
	ID          int64                `json:"id" db:"id"`
	ExternalId  string               `json:"externalId" db:"external_id"`
	UserID      int64                `json:"userId" db:"user_id"`
	Name        string               `json:"name" db:"name"`
	ChatID      string               `json:"chatId" db:"chat_id"`
	ProductName string               `json:"productName" db:"product_name"`
	Status      domain.SessionStatus `json:"status" db:"status"`
	Type        domain.SessionType   `json:"type" db:"type"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type SessionLogs struct {
	ID        int64  `json:"id" db:"id"`
	SessionID int64  `json:"sessionId" db:"session_id"`
	Text      string `json:"text" db:"text"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
