package repositories

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(ctx context.Context, name, chatId, external string) (int64, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, name, chatId, externalInfo string) (int64, error) {
	var id int64

	err := r.db.QueryRowContext(ctx, "INSERT INTO users (name, chat_id, external_info) VALUES ($1, $2, $3) RETURNING id", name, chatId, externalInfo).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
