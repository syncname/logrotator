package logrotator

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	rotator, err := NewLogRotator("./logs", time.Minute, 1024, "daily") // Интервал: 1 минута, Размер: 1 МБ
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
