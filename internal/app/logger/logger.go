package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/vo1dFl0w/users-service/internal/app/config"
)

var (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

var (
	logTimeFormat = "02/01/2006 15:04:05"
)

// Load logger with type of env from config
func LoadLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger

	options := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if v := a.Value.Any(); v != nil {
				if a.Key == slog.TimeKey {
					if t, ok := v.(time.Time); ok {
						return slog.String(slog.TimeKey, t.Format(logTimeFormat))
					}
				}
			}
			return a
		},
	}

	switch cfg.Env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo, ReplaceAttr: options.ReplaceAttr}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: options.ReplaceAttr}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug, ReplaceAttr: options.ReplaceAttr}))
	}
	
	return logger
}
