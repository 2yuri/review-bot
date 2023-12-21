package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/2yuri/review-bot/db/repositories"
	"github.com/2yuri/review-bot/domain"
	l "github.com/2yuri/review-bot/internal/logger"
	"github.com/2yuri/review-bot/internal/sqs"
	"github.com/2yuri/review-bot/internal/telegram"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strconv"
)

func Telegram(sessionRepo repositories.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var update TelegramResponse
		if err := c.ShouldBindJSON(&update); err != nil {
			log.Println("Error decoding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if update.Message.From.IsBot || update.Message.Chat.Type != "private" {
			return
		}

		if len(update.Message.Text) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message"})
			return

		}

		chatId := fmt.Sprint(update.Message.Chat.Id)

		session, err := sessionRepo.GetSessionByChatId(c, chatId)
		if err != nil {
			l.Logger.Error("Cannot find session", zap.Error(err))
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		if err := sessionRepo.AddLogs(c, session.ExternalId, update.Message.Text); err != nil {
			l.Logger.Error("Fail to add logs", zap.Error(err))
		}

		if _, err := strconv.ParseInt(string(update.Message.Text[0]), 10, 64); err != nil {
			b, err := json.Marshal(sqs.Message{
				Buttons: make([]string, 0),
				Data:    fmt.Sprint(telegram.InvalidArgMessage),
				ChatID:  chatId,
			})
			if err != nil {
				l.Logger.Error("Fail to marshal", zap.Error(err))

				c.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			}

			if err := sqs.SendToQueue(c, os.Getenv("DISPATCHER_QUEUE_URL"), string(b)); err != nil {
				l.Logger.Error("Fail to dispatch", zap.Error(err))

				c.JSON(http.StatusOK, gin.H{"status": "ok"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"status": "ok"})
			return
		}

		if err := sessionRepo.SetStatus(c, domain.Done, session.ID); err != nil {
			l.Logger.Error("Fail to set status", zap.Error(err))

			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

type TelegramResponse struct {
	UpdateId int `json:"update_id"`
	Message  struct {
		MessageId int `json:"message_id"`
		From      struct {
			Id           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			Id        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}
