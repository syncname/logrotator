package logrotator

import (
	"github.com/sirupsen/logrus"
	"log"
	"testing"
	"time"
)

func TestCustomWriter_Write(t *testing.T) {
	// Настройка LogRotator
	basePath := "./logs_logrus" // Путь к папке для логов
	interval := 24 * time.Hour
	maxSize := int64(10 << 20) // 10 MB
	strategy := DailyStrategy

	rotator, err := NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	// Создаем кастомный Writer
	writer := &LogrusWriter{rotator: rotator}

	// Настройка logrus
	logrus.SetOutput(writer)                     // Устанавливаем наш кастомный Writer
	logrus.SetFormatter(&logrus.JSONFormatter{}) // Используем JSON формат для логов
	logrus.SetLevel(logrus.DebugLevel)           // Устанавливаем уровень логирования

	// Пример использования logrus
	logrus.Info("This is an info message")
	logrus.Warn("This is a warning message")
	logrus.Error("This is an error message")

	// Закрытие текущего файла лога при завершении программы
	defer func() {
		if err := rotator.CurrentFile().Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}()
}
