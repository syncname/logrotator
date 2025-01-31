package logrotator

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
