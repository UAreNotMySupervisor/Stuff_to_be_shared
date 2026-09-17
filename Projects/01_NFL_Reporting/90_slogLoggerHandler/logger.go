package logger

import (
	"context"
	"fmt"
	"log/slog"
	repo "nfl_project_v2/02_internal/03_repository"
	"os"
)

// use this documentation https://betterstack.com/community/guides/logging/logging-in-go/

// Defining handler struct
type DbHandler struct {
	logRepo repo.LogRepo
	level   slog.Leveler
}

// Creating instance of this handler
func NewDbHandler(repository repo.LogRepo, level slog.Leveler) *DbHandler {
	return &DbHandler{
		logRepo: repository,
		level:   level,
	}
}

// Check if interface is successfully implemented
var _ slog.Handler = (*DbHandler)(nil)

// Logger_init
// Creates a new Logger based on slog.Logger but with two Hanlders. One for default Console-Logging and the second handler for database Logging.
func Logger_init(logRepo repo.LogRepo, db_log_level string) (*slog.Logger, error) {

	// All logs with level Debug and above should be logged on the console.
	// And the source (in which file) should be added
	optionsStdOut := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	// Create Json Handler based on these options and with io.Writer os.Stdout (so the console)
	handlerStdOut := slog.NewJSONHandler(os.Stdout, optionsStdOut)

	var lvl slog.Level
	// db_log_level is part of the config file and is stored there as a string
	// here the string is converted to actual slog.Levels
	switch db_log_level {
	case "Debug":
		lvl = slog.LevelDebug
	case "Info":
		lvl = slog.LevelInfo
	case "Warn":
		lvl = slog.LevelWarn
	case "Error":
		lvl = slog.LevelError
	default:
		return nil, fmt.Errorf("Error applying log level from config file. Make sure to have: 'Debug', 'Info', 'Warn' or 'Error'")
	}

	// Create a new Custome Database Handler
	databaseHandler := NewDbHandler(logRepo, lvl)

	// Create logger with two Handlers
	logger := slog.New(
		slog.NewMultiHandler(
			handlerStdOut,
			databaseHandler,
		),
	)

	return logger, nil
}

// Implementing slog.Handler
// To create a new Hanlder, the interface slog.Handler needs to be implemented (4 Functions need to be created for a struct)

func (dbH *DbHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	// only log (write to database) for this log level or higher
	if lvl >= dbH.level.Level() {
		return true
	} else {
		return false
	}
}

// With basic logging, this function can just return the original handler
func (h *DbHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

// With basic logging, this function can just return the original handler
func (h *DbHandler) WithGroup(_ string) slog.Handler {
	return h
}

// Handle
// This function is called, when there is a new log element (for example when logger.Warn() is called).
// Calls the repository function responsible for saving the entry to the database.
func (h *DbHandler) Handle(
	ctx context.Context,
	record slog.Record,
) error {
	return h.logRepo.SaveLog(ctx, record)
}
