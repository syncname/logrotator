package logrotator

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"testing"
	"time"
)

func getConsoleEncoder() zapcore.Encoder {

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func TestNewLogRotator(t *testing.T) {
	// Создаем ротатор
	rotator, err := NewLogRotator("./logs", time.Minute, 1024, DailyStrategy) // Интервал: 1 минута, Размер: 1 МБ
	if err != nil {
		fmt.Printf("Ошибка создания ротатора: %v\n", err)
		return
	}

	// Создаем zap.Logger с кастомным core
	core := NewZapCore(rotator, zapcore.InfoLevel)
	//logger := zap.New(core)

	consoleEncoder := getConsoleEncoder()
	consoleSyncer := zapcore.AddSync(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleSyncer, zapcore.InfoLevel)

	combinedCore := zapcore.NewTee(core, consoleCore)

	logger := zap.New(combinedCore, zap.AddCaller())
	defer logger.Sync() // flushes buffer, if any

	// Логируем примеры
	for i := 0; i < 1000; i++ {
		logger.Info("Пример сообщения", zap.Int("номер", i))
		time.Sleep(150 * time.Millisecond)
	}
}

func TestLogRotatorZap(t *testing.T) {
	basePath := "./logs"
	interval := 24 * time.Hour
	maxSize := int64(10 << 20)
	strategy := DailyStrategy

	rotator, err := NewLogRotator(basePath, interval, maxSize, strategy)
	if err != nil {
		log.Fatalf("Failed to create log rotator: %v", err)
	}

	zapCore := NewZapCore(rotator, zapcore.InfoLevel)
	logger := zap.New(zapCore, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	defer rotator.CurrentFile().Close()
}
