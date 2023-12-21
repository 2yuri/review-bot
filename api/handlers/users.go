package handlers

import (
	"github.com/2yuri/review-bot/db/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type createUserRequest struct {
	Name        string `json:"name"`
	ChatId      string `json:"chatId"`
	ExternalRef string `json:"externalRef"`
}

func CreateUser(userRepo repositories.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req createUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			log.Println("Error decoding JSON:", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		id, err := userRepo.CreateUser(c, req.Name, req.ChatId, req.ExternalRef)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"id": id,
		})
	}
}
