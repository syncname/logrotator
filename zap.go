package logrotator

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapCoreAdapter struct {
	rotator *LogRotator
	encoder zapcore.Encoder
	level   zapcore.Level
}

func NewZapCore(rotator *LogRotator, level zapcore.Level, encoder zapcore.Encoder) zapcore.Core {
	if encoder == nil {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}
	return &ZapCoreAdapter{
		rotator: rotator,
		encoder: encoder,
		level:   level,
	}
}

func (z *ZapCoreAdapter) Enabled(level zapcore.Level) bool {
	return level >= z.level
}

func (z *ZapCoreAdapter) With(fields []zapcore.Field) zapcore.Core {
	clone := *z
	clone.encoder = z.encoder.Clone()
	for _, field := range fields {
		field.AddTo(clone.encoder)
	}
	return &clone
}

func (z *ZapCoreAdapter) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if z.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, z)
	}
	return checkedEntry
}

func (z *ZapCoreAdapter) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	buffer, err := z.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return err
	}
	_, err = z.rotator.Write(buffer.Bytes())
	return err
}

func (z *ZapCoreAdapter) Sync() error {
	return nil
}
