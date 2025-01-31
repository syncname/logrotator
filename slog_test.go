package logrotator

import (
	"log"
	"log/slog"
	"testing"
	"time"
)

func TestSlog(t *testing.T) {
	basePath := "./logs"
	interval := 24 * time.Hour
	maxSize := int64(10 << 20)
	strategy := DailyStrategy

	rotator, err := NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	writer := &SlogWriter{}
	writer.SetRotator(rotator)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	defer rotator.CurrentFile().Close()
}
