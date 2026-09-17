package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type LogRepo struct {
	db *pgxpool.Pool
}

func NewLogRepo(db *pgxpool.Pool) *LogRepo {
	return &LogRepo{db: db}
}

type storedLog struct {
	Timestamp  time.Time      `json:"timestamp"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	Source     *slog.Source   `json:"source,omitempty"`
	Attributes map[string]any `json:"attributes,omitempty"`
}

func (l *LogRepo) SaveLog(ctx context.Context, logRecord slog.Record) error {

	attributes := make(map[string]any, logRecord.NumAttrs())

	// looping over (costume) attributes of the log records (those added with slog.Any() for example)
	logRecord.Attrs(func(attr slog.Attr) bool {
		attributes[attr.Key] = attr.Value.Any()
		return true
	})

	// Create log instance
	fullLog := storedLog{
		Timestamp:  logRecord.Time,
		Level:      logRecord.Level.String(),
		Message:    logRecord.Message,
		Source:     logRecord.Source(),
		Attributes: attributes,
	}

	// Marshal it to json (so that it can be stored in jsonb column)
	fullLogJSON, err := json.Marshal(fullLog)
	if err != nil {
		return fmt.Errorf("marshal full log: %w", err)
	}

	const query = `
		INSERT INTO logs (
			logTimestamp,
			logLevel,
			logMessage,
			fullLog
		)
		VALUES ($1, $2, $3, $4::jsonb)
	`

	// Insert log in database
	_, err = l.db.Exec(
		ctx,
		query,
		logRecord.Time,
		logRecord.Level.String(),
		logRecord.Message,
		string(fullLogJSON),
	)
	if err != nil {
		return fmt.Errorf("insert log: %w", err)
	}

	return nil
}
