package handlers

import (
	"github.com/2yuri/review-bot/db/repositories"
	"github.com/2yuri/review-bot/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type callbackRequest struct {
	Product string `json:"product"`
	UserID  int64  `json:"userId"`
	Type    string `json:"type"`
}

func Callback(sessionRepo repositories.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var review callbackRequest

		err := c.ShouldBindJSON(&review)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		if review.Product == "" || review.Type == "" || review.UserID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid request",
			})
			return
		}

		id, err := sessionRepo.CreateSession(c, review.UserID, review.Product, domain.SessionType(review.Type))
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

func GetReview(sessionRepo repositories.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status": "error",
				"error":  "invalid request",
			})

			return
		}

		review, err := sessionRepo.GetSession(c, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err.Error(),
			})
			return
		}

		c.JSON(200, review)
	}
}
