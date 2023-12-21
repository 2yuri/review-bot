package models

import "time"

type User struct {
	ID           int64  `json:"id" db:"id"`
	ChatID       string `json:"chatId" db:"chat_id"`
	Name         string `json:"name" db:"name"`
	ExternalInfo string `json:"externalInfo" db:"external_info"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
