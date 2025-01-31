package logrotator

// CustomWriter - кастомный io.Writer, который использует LogRotator
type ZeroLogWriter struct {
	rotator *LogRotator
}

// Write реализует интерфейс io.Writer для CustomWriter
func (cw *ZeroLogWriter) Write(data []byte) (int, error) {
	return cw.rotator.Write(data)
}

func (cw *ZeroLogWriter) SetRotator(rotator *LogRotator) {
	cw.rotator = rotator
}
