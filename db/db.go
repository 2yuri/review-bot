package db

import (
	"fmt"
	l "github.com/2yuri/review-bot/internal/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
)

func Start() *sqlx.DB {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connectionStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, dbName)

	db, err := sqlx.Connect("postgres", connectionStr)
	if err != nil {
		l.Logger.Fatal("Cannot start db", zap.Error(err))
	}

	return db
}
