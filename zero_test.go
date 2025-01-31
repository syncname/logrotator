package logrotator

import (
	"github.com/rs/zerolog"
	"log"
	"testing"
	"time"
)

func TestZero(t *testing.T) {
	// Настройка LogRotator
	basePath := "./logs" // Путь к папке для логов
	interval := 24 * time.Hour
	maxSize := int64(10 << 20) // 10 MB
	strategy := DailyStrategy

	rotator, err := NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	// Создаем кастомный Writer
	writer := &ZeroLogWriter{}
	writer.SetRotator(rotator)

	// Настройка zerolog
	logger := zerolog.New(writer).With().Any("loggerName", "zerolog").Timestamp().Logger()

	// Пример использования zerolog
	logger.Info().Msg("This is an info message")
	logger.Warn().Msg("This is a warning message")
	logger.Error().Msg("This is an error message")

	// Закрытие текущего файла лога при завершении программы
	defer func() {
		if err := rotator.CurrentFile().Close(); err != nil {
			log.Printf("Failed to close log file: %v", err)
		}
	}()
}
