package logrotator

type LogrusWriter struct {
	rotator *LogRotator
}

// Write реализует интерфейс io.Writer для CustomWriter
func (cw *LogrusWriter) Write(data []byte) (int, error) {
	return cw.rotator.Write(data)
}

func (cw *LogrusWriter) SetRotator(rotator *LogRotator) {
	cw.rotator = rotator
}
