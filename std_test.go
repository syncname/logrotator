package logrotator

import (
	"log"
	"testing"
	"time"
)

func TestStdLog(t *testing.T) {
	basePath := "./logs"
	interval := 24 * time.Hour
	maxSize := int64(10 << 20)
	strategy := DailyStrategy

	rotator, err := NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	writer := &StdWriter{rotator: rotator}
	log.SetOutput(writer)

	log.Println("This is a standard log message")

	defer rotator.CurrentFile().Close()
}
