package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func LoggerMiddleware(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			logger := log.With(
				"remote_addr", r.RemoteAddr,
				"http-method", r.Method,
				"path", r.URL.Path,
			)

			logger.Info("started")

			rw := &responseWriter{w, http.StatusOK}

			next.ServeHTTP(rw, r)

			var level slog.Level
			switch {
			case rw.code >= 500:
				level = slog.LevelError
			case rw.code >= 400:
				level = slog.LevelWarn
			default:
				level = slog.LevelInfo
			}

			complited := time.Since(start)
			complitedStr := fmt.Sprintf("%.3fms", float64(complited.Microseconds())/100)

			logger.Info(
				"completed",
				slog.Any(slog.LevelKey, level),
				slog.Int("code", rw.code),
				slog.String("status-text", http.StatusText(rw.code)),
				slog.String("time", complitedStr),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.code = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
