package repositories

import (
	"context"
	"github.com/2yuri/review-bot/db/models"
	"github.com/2yuri/review-bot/domain"
	"github.com/jmoiron/sqlx"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, userId int64, productName string, t domain.SessionType) (string, error)
	GetSession(ctx context.Context, externalId string) (domain.Session, error)
	GetSessionByChatId(ctx context.Context, chatId string) (domain.Session, error)
	GetPendingSessions(ctx context.Context) ([]domain.Session, error)
	SetStatus(ctx context.Context, status domain.SessionStatus, id int64) error
	AddLogs(ctx context.Context, externalId string, text string) error
}

type sessionRepository struct {
	db *sqlx.DB
}

func NewSessionRepository(db *sqlx.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) GetPendingSessions(ctx context.Context) ([]domain.Session, error) {
	var sessions []models.Session

	selectQuery := `
		WITH UniqueSessions AS (
		  SELECT
			s.user_id,
			s.id,
			s.external_id,
			s.product_name,
			u.name,
			u.chat_id,
			s.created_at,
			s.type,
			ROW_NUMBER() OVER (PARTITION BY s.user_id ORDER BY s.created_at ASC) AS row_num
		  FROM
			sessions s
			JOIN users u ON u.id = s.user_id
		  WHERE
			s.status = $1
			AND NOT EXISTS (
			  SELECT 1
			  FROM sessions
			  WHERE user_id = s.user_id AND status = $2
			)
		)
		SELECT
		  user_id,
		  id,
		  external_id,
		  product_name,
		  name,
		  chat_id,
		  created_at,
		  type
		FROM
		  UniqueSessions
		WHERE
		  row_num = 1;`

	if err := r.db.SelectContext(ctx, &sessions, selectQuery, domain.Waiting, domain.Processing); err != nil {
		return nil, err
	}

	res := make([]domain.Session, len(sessions))

	for i, session := range sessions {
		res[i].ID = session.ID
		res[i].ExternalId = session.ExternalId
		res[i].Name = session.Name
		res[i].Product = session.ProductName
		res[i].Status = session.Status
		res[i].ChatID = session.ChatID
		res[i].Type = session.Type
	}

	return res, nil
}

func (r *sessionRepository) SetStatus(ctx context.Context, status domain.SessionStatus, id int64) error {
	_, err := r.db.ExecContext(ctx, "UPDATE sessions SET status = $1 WHERE id = $2", status, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *sessionRepository) CreateSession(ctx context.Context, userId int64, productName string, t domain.SessionType) (string, error) {
	var id string

	err := r.db.QueryRowContext(ctx, "INSERT INTO sessions (user_id, product_name, type) VALUES ($1, $2, $3) RETURNING external_id", userId, productName, t).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *sessionRepository) GetSession(ctx context.Context, externalId string) (domain.Session, error) {
	var session models.Session

	selectQuery := `
		SELECT s.id, s.external_id, u.name, s.status, s.type, s.product_name
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE s.external_id = $1`

	err := r.db.QueryRowContext(ctx, selectQuery, externalId).Scan(&session.ID, &session.ExternalId, &session.Name, &session.Status, &session.Type, &session.ProductName)
	if err != nil {
		return domain.Session{}, err
	}

	var logs []models.SessionLogs

	selectLogsQuery := `
		SELECT sl.id, sl.text
		FROM session_logs sl
		WHERE sl.session_id = $1`

	if err := r.db.SelectContext(ctx, &logs, selectLogsQuery, session.ID); err != nil {
		return domain.Session{}, err
	}

	res := domain.Session{
		ID:         session.ID,
		ExternalId: session.ExternalId,
		Name:       session.Name,
		Product:    session.ProductName,
		Status:     session.Status,
		Type:       session.Type,
		Logs:       make([]domain.SessionLogs, len(logs)),
	}

	for i, log := range logs {
		res.Logs[i].ID = log.ID
		res.Logs[i].Text = log.Text
	}

	return res, err
}

func (r *sessionRepository) GetSessionByChatId(ctx context.Context, chatId string) (domain.Session, error) {
	var session models.Session

	selectQuery := `
		SELECT s.id, s.external_id, u.name, s.status, s.type, s.product_name
		FROM sessions s
		JOIN users u ON u.id = s.user_id
		WHERE u.chat_id = $1 and s.status = $2`

	err := r.db.QueryRowContext(ctx, selectQuery, chatId, domain.Processing).Scan(&session.ID, &session.ExternalId, &session.Name, &session.Status, &session.Type, &session.ProductName)
	if err != nil {
		return domain.Session{}, err
	}

	var logs []models.SessionLogs

	selectLogsQuery := `
		SELECT sl.id, sl.text
		FROM session_logs sl
		WHERE sl.session_id = $1`

	if err := r.db.SelectContext(ctx, &logs, selectLogsQuery, session.ID); err != nil {
		return domain.Session{}, err
	}

	res := domain.Session{
		ID:         session.ID,
		ExternalId: session.ExternalId,
		Name:       session.Name,
		Product:    session.ProductName,
		Status:     session.Status,
		Type:       session.Type,
		Logs:       make([]domain.SessionLogs, len(logs)),
	}

	for i, log := range logs {
		res.Logs[i].ID = log.ID
		res.Logs[i].Text = log.Text
	}

	return res, err
}

func (r *sessionRepository) AddLogs(ctx context.Context, externalId string, text string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO session_logs (session_id, text) VALUES ((SELECT id FROM sessions WHERE external_id = $1), $2)", externalId, text)
	if err != nil {
		return err
	}

	return nil
}
