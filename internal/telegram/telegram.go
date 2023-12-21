package telegram

import (
	"encoding/json"
	"fmt"
	"github.com/2yuri/review-bot/internal/env"
	"github.com/2yuri/review-bot/internal/http"
	l "github.com/2yuri/review-bot/internal/logger"
	"go.uber.org/zap"
)

const ReviewMessage = "We noticed you've recently received your %s. We'd love to hear about your experience. \nCan you rate it 1-5 on how you like it?"
const InvalidArgMessage = "Invalid Option, please select a valid button. \nCan you rate it 1-5 on how you like it?"

func SendBotMessage(chatId, text string, buttons []string) error {
	req := map[string]interface{}{
		"chat_id": chatId,
		"text":    text,
	}

	if len(buttons) != 0 {
		btns := make([]map[string]string, len(buttons))
		for i, btn := range buttons {
			btns[i] = map[string]string{
				"text":          btn,
				"callback_data": "review-" + fmt.Sprint(i),
			}
		}

		keyboardOpts := map[string]interface{}{
			"keyboard": [][]map[string]string{
				btns,
			},
			"one_time_keyboard": true,
		}

		kBytes, _ := json.Marshal(keyboardOpts)

		req["reply_markup"] = string(kBytes)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", env.Get("TELEGRAM_BOT_KEY").String())

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	l.Logger.Debug("Request", zap.Any("req", req))

	return http.Post(url, headers, req, nil)
}
