package api

import (
	"github.com/2yuri/review-bot/api/handlers"
	"github.com/2yuri/review-bot/db/repositories"
	"github.com/gin-gonic/gin"
)

type api struct {
	port string

	gin         *gin.Engine
	userRepo    repositories.UserRepository
	sessionRepo repositories.SessionRepository
}

func New(port string, uR repositories.UserRepository, sR repositories.SessionRepository) *api {
	eng := gin.Default()

	return &api{port: port, gin: eng, userRepo: uR, sessionRepo: sR}
}

func (a *api) setupRoutes() {
	a.gin.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	group := a.gin.Group("/v1")

	group.POST("/callback", handlers.Callback(a.sessionRepo))
	group.GET("/review/:id", handlers.GetReview(a.sessionRepo))
	group.POST("/user", handlers.CreateUser(a.userRepo))
	group.POST("/telegram", handlers.Telegram(a.sessionRepo))
}

func (a *api) Start() {
	a.setupRoutes()

	a.gin.Run(":" + a.port)
}
