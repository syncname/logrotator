package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogRotator struct {
	mu               sync.Mutex
	basePath         string
	interval         time.Duration
	maxSize          int64
	currentFile      *os.File
	currentFileSize  int64
	rotationStrategy string
	lastRotationTime time.Time
}

// NewLogRotator создает новый экземпляр LogRotator
func NewLogRotator(basePath string, interval time.Duration, maxSize int64, strategy string) (*LogRotator, error) {
	rotator := &LogRotator{
		basePath:         basePath,
		interval:         interval,
		maxSize:          maxSize,
		rotationStrategy: strategy,
		lastRotationTime: time.Now(),
	}

	//create new log file
	if err := rotator.rotate(); err != nil {
		return nil, err
	}

	return rotator, nil
}

// Write data to log file
func (r *LogRotator) Write(data []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	//time log rotation
	if time.Since(r.lastRotationTime) >= r.interval {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	//size log rotation
	if r.currentFileSize+int64(len(data)) > r.maxSize {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	n, err := r.currentFile.Write(data)
	if err != nil {
		return 0, err
	}
	r.currentFileSize += int64(n)

	return n, nil
}

func (r *LogRotator) rotate() error {

	if r.currentFile != nil {
		if err := r.currentFile.Close(); err != nil {
			return err
		}
	}

	folder := r.getRotationFolder()
	if err := os.MkdirAll(folder, 0755); err != nil {
		return fmt.Errorf("mkdir error: %w", err)
	}

	filename := filepath.Join(folder, fmt.Sprintf("log_%d.log", time.Now().Unix()))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("create log file error: %w", err)
	}

	r.currentFile = file
	r.currentFileSize = 0
	r.lastRotationTime = time.Now()

	return nil
}

// getRotationFolder возвращает путь для текущей ротации в зависимости от стратегии
func (r *LogRotator) getRotationFolder() string {
	now := time.Now()
	switch r.rotationStrategy {
	case "daily":
		return filepath.Join(r.basePath, now.Format("2006-01-02"))
	case "weekly":
		year, week := now.ISOWeek()
		return filepath.Join(r.basePath, fmt.Sprintf("%d-W%02d", year, week))
	case "monthly":
		return filepath.Join(r.basePath, now.Format("2006-01"))
	case "yearly":
		return filepath.Join(r.basePath, now.Format("2006"))
	default:
		return r.basePath
	}
}

// ZapCoreAdapter адаптирует LogRotator для использования с zapcore.Core
type ZapCoreAdapter struct {
	rotator *LogRotator
	encoder zapcore.Encoder
	level   zapcore.Level
}

// NewZapCore создает новый zapcore.Core с использованием кастомного ротатора
func NewZapCore(rotator *LogRotator, level zapcore.Level) zapcore.Core {
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	return &ZapCoreAdapter{
		rotator: rotator,
		encoder: encoder,
		level:   level,
	}
}

// Enabled проверяет, активен ли текущий уровень логирования
func (z *ZapCoreAdapter) Enabled(level zapcore.Level) bool {
	return level >= z.level
}

// With добавляет дополнительные поля в лог
func (z *ZapCoreAdapter) With(fields []zapcore.Field) zapcore.Core {
	clone := *z
	clone.encoder = z.encoder.Clone()
	for _, field := range fields {
		field.AddTo(clone.encoder)
	}
	return &clone
}

// Check проверяет, можно ли логировать сообщение
func (z *ZapCoreAdapter) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, z)
	}
	return checkedEntry
}

// Write пишет сообщение в лог
func (z *ZapCoreAdapter) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buffer, err := z.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}
	_, err = z.rotator.Write(buffer.Bytes())
	return err
}

// Sync выполняет синхронизацию (для совместимости с zapcore.Core)
func (z *ZapCoreAdapter) Sync() error {
	return nil
}

func main() {
	// Создаем ротатор
	rotator, err := NewLogRotator("./logs", time.Minute, 1024*1024, "daily") // Интервал: 1 минута, Размер: 1 МБ
	if err != nil {
		fmt.Printf("Ошибка создания ротатора: %v\n", err)
		return
	}

	// Создаем zap.Logger с кастомным core
	core := NewZapCore(rotator, zapcore.InfoLevel)
	logger := zap.New(core)

	// Логируем примеры
	for i := 0; i < 1000; i++ {
		logger.Info("Пример сообщения", zap.Int("номер", i))
		time.Sleep(500 * time.Millisecond)
	}
}
