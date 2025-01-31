package logrotator

import (
	"github.com/sirupsen/logrus"
	"log"
	"testing"
	"time"
)

func TestCustomWriter_Write(t *testing.T) {
	basePath := "./logs_logrus"
	interval := 24 * time.Hour
	maxSize := int64(10 << 20) // 10 MB
	strategy := DailyStrategy

	rotator, err := NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}
	writer := &LogrusWriter{rotator: rotator}

	logrus.SetOutput(writer)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	logrus.Info("This is an info message")
	logrus.Warn("This is a warning message")
	logrus.Error("This is an error message")

	defer func() {
		if err := rotator.CurrentFile().Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}()
}
